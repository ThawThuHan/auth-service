package config

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

func NewLogger(env string) *zap.Logger {
	var logConfig zap.Config

	if env == "production" {
		println("Running in production mode")
		// Ensure logs directory exists
		logsDir := filepath.Join(".", "logs")
		if err := os.MkdirAll(logsDir, 0755); err != nil {
			panic("Failed to create logs directory: " + err.Error())
		}
		logConfig = zap.NewProductionConfig()
		logConfig.OutputPaths = []string{"stdout", filepath.Join(logsDir, "app.log")}
		logConfig.EncoderConfig.TimeKey = "timestamp"
	} else {
		logConfig = zap.NewDevelopmentConfig()
		logConfig.EncoderConfig.TimeKey = "ts"
	}

	logger, err := logConfig.Build()
	if err != nil {
		panic("Failed to create logger: " + err.Error())
	}
	return logger
}
