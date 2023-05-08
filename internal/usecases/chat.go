//go:generate mockgen -source=$GOFILE -destination=../tests/mocks/$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package usecases

import (
	"context"

	"github.com/swallowarc/simple-line-ai-bot/internal/domain"
)

type (
	Chat interface {
		Help(ctx context.Context, callback Callback) error
		ClearChatHistory(ctx context.Context, es domain.EventSource, callback Callback) error
		Chat(ctx context.Context, es domain.EventSource, req string, callback Callback) error
	}

	chat struct {
		cacheRepo  CacheRepository
		openAIRepo OpenAIRepository
	}

	Callback func(ctx context.Context, replyMessage string) error
)

func NewChat(
	memDBRepo CacheRepository,
	openAIRepo OpenAIRepository,
) Chat {
	return &chat{
		cacheRepo:  memDBRepo,
		openAIRepo: openAIRepo,
	}
}

func (uc *chat) Help(ctx context.Context, callback Callback) error {
	return callback(ctx, domain.MessageHelp())
}

func (uc *chat) ClearChatHistory(ctx context.Context, es domain.EventSource, callback Callback) error {
	if err := uc.cacheRepo.DeleteChatMessages(ctx, es); err != nil {
		return err
	}

	if err := callback(ctx, domain.MessageClearHistory); err != nil {
		return err
	}

	return nil
}

func (uc *chat) Chat(ctx context.Context, es domain.EventSource, req string, callback Callback) error {
	// get chat history
	messages, err := uc.cacheRepo.ListChatMessages(ctx, es)
	if err != nil {
		return err
	}
	messages = append(messages, domain.ChatMessage{
		Role:    domain.RoleUser,
		Content: req,
	})

	// request to openAI
	res, err := uc.openAIRepo.ChatCompletion(ctx, messages)
	if err != nil {
		return err
	}
	messages = append(messages, res...)
	if l := len(messages); l > domain.ChatHistoryLimit {
		messages = messages[l-domain.ChatHistoryLimit:]
	}

	if err := uc.cacheRepo.SetChatMessages(ctx, es, messages); err != nil {
		return err
	}

	lm := res.LatestMessage()
	if lm != nil {
		return callback(ctx, lm.Content)
	}

	return nil
}
