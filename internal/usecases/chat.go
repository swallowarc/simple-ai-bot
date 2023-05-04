package usecases

import (
	"context"

	"github.com/swallowarc/simple-line-ai-bot/internal/domain"
)

type (
	Chat struct {
		cacheRepo  CacheRepository
		openAIRepo OpenAIRepository
	}
)

func NewChat(
	memDBRepo CacheRepository,
	openAIRepo OpenAIRepository,
) *Chat {
	return &Chat{
		cacheRepo:  memDBRepo,
		openAIRepo: openAIRepo,
	}
}

func (uc *Chat) Chat(ctx context.Context, es domain.EventSource, req string) (string, error) {
	messages, err := uc.cacheRepo.ListChatMessages(ctx, es)
	if err != nil {
		return "", err
	}

	messages = append(messages, domain.ChatMessage{
		Role:    domain.RoleUser,
		Content: req,
	})

	newMessages, err := uc.openAIRepo.ChatCompletion(ctx, messages)
	if err != nil {
		return "", err
	}

	messages = append(messages, newMessages...)

	// TODO: 最新の10件をmemDに保存

	if lm := newMessages.LatestMessage(); lm != nil {
		return lm.Content, nil
	}

	return "", nil
}
