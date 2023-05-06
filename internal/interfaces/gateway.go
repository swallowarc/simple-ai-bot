//go:generate mockgen -source=$GOFILE -destination=../tests/mocks/$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package interfaces

import (
	"context"
	"time"
)

type MemDBClient interface {
	Ping(ctx context.Context) error
	Expire(ctx context.Context, key string, duration time.Duration) error
	Set(ctx context.Context, key string, value any, duration time.Duration) error
	SetNX(ctx context.Context, key string, value any, duration time.Duration) (bool, error)
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}
