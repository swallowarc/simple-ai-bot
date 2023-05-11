package repositories

import (
	"context"
	"fmt"

	"github.com/swallowarc/simple-line-ai-bot/internal/domain"
	"github.com/swallowarc/simple-line-ai-bot/internal/usecases"

	"github.com/swallowarc/simple-line-ai-bot/internal/interfaces"
)

type (
	chatRepository struct {
		memDBCli interfaces.MemDBClient
		aiCli    interfaces.AIClient
	}
)

const (
	keyChatMessages = "chat_messages:%s"
)

func NewChatRepository(
	memDBCli interfaces.MemDBClient,
	aiCli interfaces.AIClient,
) usecases.ChatRepository {
	return &chatRepository{
		memDBCli: memDBCli,
		aiCli:    aiCli,
	}
}

func (r *chatRepository) chatMessagesKey(es domain.EventSource) string {
	return fmt.Sprintf(keyChatMessages, es.UniqueID())
}

func (r *chatRepository) ListCacheMessages(ctx context.Context, es domain.EventSource) (domain.ChatMessages, error) {
	return getFromMemDB[domain.ChatMessages](ctx, r.memDBCli, r.chatMessagesKey(es))
}

func (r *chatRepository) UpsertCacheMessages(ctx context.Context, es domain.EventSource, cms domain.ChatMessages) error {
	return setToMemDB[domain.ChatMessages](ctx, r.memDBCli, r.chatMessagesKey(es), cms, domain.ChatHistoryTTL)
}

func (r *chatRepository) DeleteCacheMessages(ctx context.Context, es domain.EventSource) error {
	return r.memDBCli.Del(ctx, r.chatMessagesKey(es))
}

func (r *chatRepository) Chat(ctx context.Context, messages domain.ChatMessages) (domain.ChatMessages, error) {
	return r.aiCli.ChatCompletion(ctx, messages)
}
