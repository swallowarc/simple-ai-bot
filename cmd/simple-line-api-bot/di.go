package main

import (
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/swallowarc/lime/lime"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/swallowarc/simple-line-ai-bot/internal/core"
	"github.com/swallowarc/simple-line-ai-bot/internal/infrastructures/env"
	"github.com/swallowarc/simple-line-ai-bot/internal/infrastructures/redis"
	"github.com/swallowarc/simple-line-ai-bot/internal/interfaces"
	"github.com/swallowarc/simple-line-ai-bot/internal/interfaces/repositories"
	"github.com/swallowarc/simple-line-ai-bot/internal/usecases"
)

func coreModules() fx.Option {
	return fx.Module(""+
		"core",
		fx.Provide(
			zap.NewProductionConfig,
			core.NewZapLoggerWithConfig,
		),
	)
}

func infrastructureModules() fx.Option {
	return fx.Module("infrastructures",
		fx.Provide(
			env.GetBotEnv,
			env.GetRedisEnv,
			env.GetLimeEnv,
			redis.NewClient,
			func() *http.Client {
				return http.DefaultClient
			},
			fx.Annotate(
				func(h lime.EventHandler) lime.APIServerOption {
					return lime.WithEventHandler(linebot.EventTypeMessage, h)
				},
				fx.ResultTags(`group:"lime_options"`),
			),
			fx.Annotate(
				func(l *zap.Logger) lime.APIServerOption {
					return lime.WithLogger(l)
				},
				fx.ResultTags(`group:"lime_options"`),
			),
			func(p struct {
				fx.In
				LimeEnv     lime.Env
				LimeOptions []lime.APIServerOption `group:"lime_options"`
			}) (lime.APIServer, error) {
				return lime.NewServer(p.LimeEnv, p.LimeOptions...)
			},
		),
	)
}

func interfaceModules() fx.Option {
	return fx.Module("interfaces",
		fx.Provide(
			interfaces.NewMessageEventHandler,
			repositories.NewCacheRepository,
			repositories.NewOpenAIRepository,
		),
	)
}

func usecaseModules() fx.Option {
	return fx.Module("usecases",
		fx.Provide(
			usecases.NewChat,
		),
	)
}
