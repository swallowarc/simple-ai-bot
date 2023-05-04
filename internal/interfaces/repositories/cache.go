package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/swallowarc/simple-line-ai-bot/internal/core"
	"github.com/swallowarc/simple-line-ai-bot/internal/domain"
	"github.com/swallowarc/simple-line-ai-bot/internal/usecases"

	"github.com/swallowarc/simple-line-ai-bot/internal/interfaces"
)

type (
	cacheRepository struct {
		memDBCli interfaces.MemDBClient
	}
)

const (
	keyChatMessages = "chat_messages:%s:%s"
)

func NewCacheRepository(memDBCli interfaces.MemDBClient) usecases.CacheRepository {
	return &cacheRepository{
		memDBCli: memDBCli,
	}
}

func (c *cacheRepository) ListChatMessages(ctx context.Context, es domain.EventSource) (domain.ChatMessages, error) {
	j, err := c.memDBCli.Get(ctx, fmt.Sprintf(keyChatMessages, es.Type, es.ID))
	if err != nil {
		if !errors.Is(err, core.ErrNotFound) {
			return domain.ChatMessages{}, err
		}
		return nil, err
	}

	var cms domain.ChatMessages
	if err := json.Unmarshal([]byte(j), &cms); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal json")
	}

	return domain.ChatMessages{}, nil
}