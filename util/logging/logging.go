package logging

import "go.uber.org/zap"

var logger *zap.Logger

type LoggerConfig struct{}

func GetDefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{}
}

func GetLogger(config *LoggerConfig) *zap.Logger {
	if logger != nil {
		return logger
	}
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return logger
}
