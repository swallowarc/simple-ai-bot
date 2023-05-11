package env

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/swallowarc/lime/lime"

	"github.com/swallowarc/simple-line-ai-bot/internal/infrastructures/line"
	"github.com/swallowarc/simple-line-ai-bot/internal/infrastructures/openai"
	"github.com/swallowarc/simple-line-ai-bot/internal/infrastructures/redis"
)

type Env struct {
	Env string `envconfig:"env" default:"DEBUG"`

	// LicenseMode is a flag that indicates whether to use the license mode.
	//   If true, the administrator's approval is required in order to use the bot's AI.
	LicenseMode bool `envconfig:"license_mode" default:"false"`
	// AdminLineUserID is the user ID of the LINE account that will be used to manage the bot.
	//   Must be specified if LicenseMode=true.
	//   This user ID can be obtained from the LINE Developers Console.
	//     https://developers.line.biz/console/channel/{channel_id}}/basics
	AdminLineUserID string `envconfig:"admin_user_id" default:"mock-admin-user-id"`
}

var (
	botEnv    Env
	redisEnv  redis.Env
	limeEnv   lime.Env
	lineEnv   line.Env
	openAIEnv openai.Env
)

func init() {
	envconfig.MustProcess("", &botEnv)
	envconfig.MustProcess("redis", &redisEnv)
	envconfig.MustProcess("lime", &limeEnv)
	envconfig.MustProcess("line", &lineEnv)
	envconfig.MustProcess("open_ai", &openAIEnv)
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

func GetLineEnv() line.Env {
	return lineEnv
}

func GetOpenAI() openai.Env {
	return openAIEnv
}
