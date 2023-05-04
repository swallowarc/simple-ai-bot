package usecases

import (
	"context"

	"github.com/swallowarc/simple-line-ai-bot/internal/domain"
)

type (
	CacheRepository interface {
		ListChatMessages(ctx context.Context, es domain.EventSource) (domain.ChatMessages, error)
		SetChatMessages(ctx context.Context, es domain.EventSource, cms domain.ChatMessages) error
	}
	OpenAIRepository interface {
		ChatCompletion(ctx context.Context, messages domain.ChatMessages) (domain.ChatMessages, error)
	}
)
