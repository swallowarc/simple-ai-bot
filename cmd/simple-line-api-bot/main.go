package main

import (
	"context"
	"log"

	"github.com/swallowarc/lime/lime"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/swallowarc/simple-line-ai-bot/internal/infrastructures/env"
	"github.com/swallowarc/simple-line-ai-bot/internal/infrastructures/redis"
	"github.com/swallowarc/simple-line-ai-bot/internal/interfaces"
)

func main() {
	app := fx.New(
		fx.WithLogger(func(log *zap.Logger, env env.BotEnv) fxevent.Logger {
			if env.Env == "DEBUG" {
				return &fxevent.ZapLogger{Logger: log}
			}
			return fxevent.NopLogger
		}),
		coreModules(),
		infrastructureModules(),
		interfaceModules(),
		usecaseModules(),
		fx.Invoke(initialize, run),
	)
	if err := app.Err(); err != nil {
		log.Fatalf(err.Error())
	}

	app.Run()
}

func initialize(lc fx.Lifecycle, logger *zap.Logger, cli interfaces.MemDBClient, conf redis.Config) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Info("start")
			if err := cli.Ping(context.Background()); err != nil {
				logger.Fatal("failed to ping to redis", zap.Error(err), zap.String("redis_host_port", conf.HostPort))
				return err
			}
			logger.Info("ping to redis was successful")
			return nil
		},
		OnStop: func(context.Context) error {
			logger.Info("stop")
			return nil
		},
	})
}

func run(s lime.APIServer, logger *zap.Logger) {
	go func() {
		if err := s.Start(); err != nil {
			logger.Error("failed to start api server", zap.Error(err))
		}
	}()
}
