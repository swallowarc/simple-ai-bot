package repositories

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

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

func (c *cacheRepository) chatMessagesKey(es domain.EventSource) string {
	return fmt.Sprintf(keyChatMessages, es.Type, es.ID)
}

func (c *cacheRepository) ListChatMessages(ctx context.Context, es domain.EventSource) (domain.ChatMessages, error) {
	j, err := c.memDBCli.Get(ctx, c.chatMessagesKey(es))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ChatMessages{}, nil
		}
		return nil, err
	}

	var cms domain.ChatMessages
	if err := json.Unmarshal([]byte(j), &cms); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal json")
	}

	return domain.ChatMessages{}, nil
}

func (c *cacheRepository) SetChatMessages(ctx context.Context, es domain.EventSource, cms domain.ChatMessages) error {
	j, err := json.Marshal(cms)
	if err != nil {
		return errors.Wrap(err, "failed to marshal json")
	}

	if err := c.memDBCli.Set(ctx, c.chatMessagesKey(es), string(j), domain.ChatHistoryLife); err != nil {
		return err
	}

	return nil
}
