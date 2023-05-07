package env

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/swallowarc/lime/lime"

	"github.com/swallowarc/simple-line-ai-bot/internal/infrastructures/redis"
)

type (
	Env struct {
		Env                  string  `envconfig:"env" default:"DEBUG"`
		OpenAIAPIKey         string  `envconfig:"open_ai_api_key" required:"true"`
		OpenAIAPIMaxTokens   int     `envconfig:"open_ai_api_max_tokens" default:"400"`
		OpenAIAPITemperature float64 `envconfig:"open_ai_api_temperature" default:"0.6"`
	}
)

var (
	botEnv   Env
	redisEnv redis.Env
	limeEnv  lime.Env
)

func init() {
	envconfig.MustProcess("", &botEnv)
	envconfig.MustProcess("redis", &redisEnv)
	envconfig.MustProcess("lime", &limeEnv)
}

func GetBotEnv() Env {
	return botEnv
}

func GetRedisEnv() redis.Env {
	return redisEnv
}

func GetLimeEnv() lime.Env {
	return limeEnv
}
