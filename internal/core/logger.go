package core

import (
	"log"

	"go.uber.org/zap"
)

func NewZapLoggerWithConfig(config zap.Config) *zap.Logger {
	return build(config)
}

// build get zap logger instance.
func build(config zap.Config) *zap.Logger {
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	logger, err := config.Build()
	if err != nil {
		log.Println("create logger failed")
		log.Fatal(err)
	}
	return logger
}
