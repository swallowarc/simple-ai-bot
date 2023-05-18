//go:generate mockgen -source=$GOFILE -destination=../tests/mocks/$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package usecases

import (
	"context"

	"github.com/swallowarc/simple-line-ai-bot/internal/domain"
)

type (
	Chat interface {
		Help(ctx context.Context, replyToken string) error
		ClearChatHistory(ctx context.Context, es domain.EventSource, replyToken string) error
		Chat(ctx context.Context, es domain.EventSource, replyToken, req string) error
	}

	chat struct {
		chatRepo ChatRepository
		msgRepo  MessagingRepository
	}
)

func NewChat(
	chatRepo ChatRepository,
	msgRepo MessagingRepository,
) Chat {
	return &chat{
		chatRepo: chatRepo,
		msgRepo:  msgRepo,
	}
}

func (u *chat) Help(ctx context.Context, replyToken string) error {
	if err := u.msgRepo.ReplyMessages(ctx, replyToken, domain.MessageHelp()); err != nil {
		return err
	}

	return nil
}

func (u *chat) ClearChatHistory(ctx context.Context, es domain.EventSource, replyToken string) error {
	if err := u.chatRepo.DeleteCacheMessages(ctx, es); err != nil {
		return err
	}

	if err := u.msgRepo.ReplyMessages(ctx, replyToken, domain.MessageClearHistory); err != nil {
		return err
	}

	return nil
}

func (u *chat) Chat(ctx context.Context, es domain.EventSource, replyToken, req string) error {
	// get chat history
	messages, err := u.chatRepo.ListCacheMessages(ctx, es)
	if err != nil {
		return err
	}
	messages = append(messages, domain.ChatMessage{
		Role:    domain.RoleUser,
		Content: req,
	})

	// request to openAI
	res, err := u.chatRepo.Chat(ctx, messages)
	if err != nil {
		return err
	}
	messages = append(messages, res...)
	if l := len(messages); l > domain.ChatHistoryLimit {
		messages = messages[l-domain.ChatHistoryLimit:]
	}

	if err := u.chatRepo.UpsertCacheMessages(ctx, es, messages); err != nil {
		return err
	}

	// reply latest message
	if lm := res.LatestMessage(); lm != nil {
		if err := u.msgRepo.ReplyMessages(ctx, replyToken, lm.Content); err != nil {
			return err
		}
	}

	return nil
}
