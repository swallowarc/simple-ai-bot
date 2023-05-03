package usecases

import (
	"context"

	"github.com/swallowarc/simple-line-ai-bot/internal/domain"
	"github.com/swallowarc/simple-line-ai-bot/internal/interfaces"
)

type (
	Chat struct {
		memDBRepo  interfaces.MemDBClient
		openAIRepo OpenAIRepository
	}
)

const (
	chatHistoryLimit = 10
)

func NewChat(
	memDBRepo interfaces.MemDBClient,
	openAIRepo OpenAIRepository,
) *Chat {
	return &Chat{
		memDBRepo:  memDBRepo,
		openAIRepo: openAIRepo,
	}
}

func (c *Chat) Chat(ctx context.Context, es domain.EventSource, req string) (string, error) {
	c.memDBRepo.Get(ctx, es.Key())
	return nil
}
