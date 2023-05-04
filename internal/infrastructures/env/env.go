package env

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/swallowarc/lime/lime"

	"github.com/swallowarc/simple-line-ai-bot/internal/infrastructures/redis"
)

type (
	BotEnv struct {
		Env          string `envconfig:"env" default:"DEBUG"`
		OpenAIAPIKey string `envconfig:"open_ai_api_key" required:"true"`
	}
)

var (
	Bot   BotEnv
	Redis redis.Config
	Lime  lime.Env
)

func init() {
	envconfig.MustProcess("", &Bot)
	envconfig.MustProcess("redis", &Redis)
	envconfig.MustProcess("lime", &Lime)
}

func GetBotEnv() BotEnv {
	return Bot
}

func GetRedisEnv() redis.Config {
	return Redis
}

func GetLimeEnv() lime.Env {
	return Lime
}
