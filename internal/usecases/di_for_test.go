package usecases_test

import (
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/swallowarc/simple-line-ai-bot/internal/infrastructures/env"
	"github.com/swallowarc/simple-line-ai-bot/internal/infrastructures/redis"
	"github.com/swallowarc/simple-line-ai-bot/internal/interfaces/repositories"
	"github.com/swallowarc/simple-line-ai-bot/internal/usecases"
)

func coreModules() fx.Option {
	return fx.Options(
		fx.Provide(
			zap.NewExample,
		),
	)
}

func infrastructureModules() fx.Option {
	return fx.Options(
		fx.Provide(
			env.GetBotEnv,
			env.GetRedisEnv,
			redis.NewClient,
			func() *http.Client {
				return http.DefaultClient
			},
		),
	)
}

func interfaceModules() fx.Option {
	return fx.Options(
		fx.Provide(
			repositories.NewCacheRepository,
			func(cli *http.Client, e env.Env) usecases.OpenAIRepository {
				return repositories.NewOpenAIRepository(cli, e.OpenAIAPIKey, e.OpenAIAPIMaxTokens, e.OpenAIAPITemperature)
			},
		),
	)
}

func usecaseModules() fx.Option {
	return fx.Options(
		fx.Provide(
			usecases.NewChat,
		),
	)
}
