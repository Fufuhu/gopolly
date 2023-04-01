package logging

import "go.uber.org/zap"

var logger *zap.Logger

func GetLogger() *zap.Logger {
	if logger != nil {
		return logger
	}
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return logger
}
