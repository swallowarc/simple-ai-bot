package main

import (
	"net/http"

	"github.com/swallowarc/lime/lime"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/swallowarc/simple-line-ai-bot/internal/core"
	"github.com/swallowarc/simple-line-ai-bot/internal/infrastructures/env"
	"github.com/swallowarc/simple-line-ai-bot/internal/infrastructures/line"
	"github.com/swallowarc/simple-line-ai-bot/internal/infrastructures/openai"
	"github.com/swallowarc/simple-line-ai-bot/internal/infrastructures/redis"
	"github.com/swallowarc/simple-line-ai-bot/internal/interfaces/eventhandler"
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
			// env
			env.GetBotEnv,
			env.GetRedisEnv,
			env.GetLimeEnv,
			env.GetLineEnv,
			env.GetOpenAI,

			// clients
			redis.NewClient,
			func() *http.Client {
				// TODO: replace to retryable client.
				return http.DefaultClient
			},
			line.NewClient,
			openai.NewClient,

			// lime options
			fx.Annotate(
				func(h *eventhandler.Follow) lime.APIServerOption {
					return lime.WithEventHandler(h)
				},
				fx.ResultTags(`group:"lime_options"`),
			),
			fx.Annotate(
				func(h *eventhandler.Join) lime.APIServerOption {
					return lime.WithEventHandler(h)
				},
				fx.ResultTags(`group:"lime_options"`),
			),
			fx.Annotate(
				func(h *eventhandler.Leave) lime.APIServerOption {
					return lime.WithEventHandler(h)
				},
				fx.ResultTags(`group:"lime_options"`),
			),
			fx.Annotate(
				func(h *eventhandler.Message) lime.APIServerOption {
					return lime.WithEventHandler(h)
				},
				fx.ResultTags(`group:"lime_options"`),
			),
			fx.Annotate(
				func(h *eventhandler.Unfollow) lime.APIServerOption {
					return lime.WithEventHandler(h)
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
			}) lime.APIServer {
				return lime.NewServer(p.LimeEnv, p.LimeOptions...)
			},
		),
	)
}

func interfaceModules() fx.Option {
	return fx.Module("interfaces",
		fx.Provide(
			// event handlers
			eventhandler.NewMessageEventHandler,
			eventhandler.NewJoinEventHandler,
			eventhandler.NewLeaveHandler,
			eventhandler.NewFollowHandler,
			eventhandler.NewUnfollowHandler,
			// repositories
			repositories.NewChatRepository,
			repositories.NewMessagingRepository,
			repositories.NewLicenseRepository,
		),
	)
}

func usecaseModules() fx.Option {
	return fx.Module("usecases",
		fx.Provide(
			usecases.NewChat,
			func(
				msgRepo usecases.MessagingRepository,
				lineRepo usecases.LicenseRepository,
				env env.Env,
			) usecases.License {
				return usecases.NewLicense(msgRepo, lineRepo, env.LicenseMode, env.AdminLineUserID)
			},
		),
	)
}
