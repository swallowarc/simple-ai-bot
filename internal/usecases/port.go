//go:generate mockgen -source=$GOFILE -destination=../tests/mocks/$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package usecases

import (
	"context"

	"github.com/swallowarc/simple-line-ai-bot/internal/domain"
)

type (
	CacheRepository interface {
		ListChatMessages(ctx context.Context, es domain.EventSource) (domain.ChatMessages, error)
		SetChatMessages(ctx context.Context, es domain.EventSource, cms domain.ChatMessages) error
		DeleteChatMessages(ctx context.Context, es domain.EventSource) error
	}
	OpenAIRepository interface {
		ChatCompletion(ctx context.Context, messages domain.ChatMessages) (domain.ChatMessages, error)
	}
)
