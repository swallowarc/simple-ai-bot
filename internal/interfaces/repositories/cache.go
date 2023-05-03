package repositories

import (
	"context"

	"github.com/swallowarc/simple-line-ai-bot/internal/interfaces"
)

type (
	cacheRepository struct {
		memDBCli interfaces.MemDBClient
	}
)

const (
	keyChatHistory = "ch:%s"
)

func NewCacheRepository(memDBCli interfaces.MemDBClient) *cacheRepository {
	return &cacheRepository{
		memDBCli: memDBCli,
	}
}

func (c *cacheRepository) Get(ctx context.Context, key string) (string, error) {
	return c.memDBCli.Get(ctx, key)
}
