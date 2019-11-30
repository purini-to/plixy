package log

import (
	"fmt"

	"go.uber.org/zap"
)

var w *writer

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("Could not initial build logger err: %v", err))
	}

	w = &writer{logger: logger}
}

type writer struct {
	logger *zap.Logger
}

// SetLogger is set log writer
func SetLogger(logger *zap.Logger) {
	w = &writer{logger: logger}
}

// GetLogger is get log writer
func GetLogger() *zap.Logger {
	return w.logger
}

func Debug(msg string, fields ...zap.Field) {
	w.logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	w.logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	w.logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	w.logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	w.logger.Fatal(msg, fields...)
}
