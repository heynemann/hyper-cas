package utils

import (
	"go.uber.org/zap"
)

var instance *zap.Logger

func LoggerInstance() *zap.Logger {
	var err error
	if instance == nil {
		instance, err = zap.NewProduction()
		if err != nil {
			panic(err)
		}
	}
	return instance
}

func SetDebug() {
	var err error
	instance, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
}

func LogDebug(message string, fields ...zap.Field) {
	logger := LoggerInstance()
	defer logger.Sync()
	logger.Debug(message, fields...)
}

func LogInfo(message string, fields ...zap.Field) {
	logger := LoggerInstance()
	defer logger.Sync()
	logger.Info(message, fields...)
}

func LogWarn(message string, fields ...zap.Field) {
	logger := LoggerInstance()
	defer logger.Sync()
	logger.Warn(message, fields...)
}

func LogError(message string, fields ...zap.Field) {
	logger := LoggerInstance()
	defer logger.Sync()
	logger.Error(message, fields...)
}

func LogFatal(message string, fields ...zap.Field) {
	logger := LoggerInstance()
	defer logger.Sync()
	logger.Fatal(message, fields...)
}

func LoggerWith(fields ...zap.Field) *zap.Logger {
	logger := LoggerInstance()
	return logger.With(fields...)
}
