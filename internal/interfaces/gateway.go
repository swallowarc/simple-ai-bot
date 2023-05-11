//go:generate mockgen -source=$GOFILE -destination=../tests/mocks/$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package interfaces

import (
	"context"
	"time"

	"github.com/swallowarc/simple-line-ai-bot/internal/domain"
)

type (
	MemDBClient interface {
		Ping(ctx context.Context) error
		Expire(ctx context.Context, key string, duration time.Duration) error
		Set(ctx context.Context, key string, value any, duration time.Duration) error
		SetNX(ctx context.Context, key string, value any, duration time.Duration) (bool, error)
		SetXX(ctx context.Context, key string, value any, duration time.Duration) (bool, error)
		Get(ctx context.Context, key string) (string, error)
		Del(ctx context.Context, key string) error
	}

	AIClient interface {
		ChatCompletion(ctx context.Context, messages domain.ChatMessages) (domain.ChatMessages, error)
	}
)
