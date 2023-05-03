package usecases

import (
	"context"

	"github.com/swallowarc/simple-line-ai-bot/internal/interfaces/repositories"
)

type OpenAIRepository interface {
	ChatCompletion(ctx context.Context, req *repositories.ChatCompletionRequest) (*repositories.ChatCompletionResponse, error)
}
