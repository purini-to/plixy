package log

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type loggerKeyType int

const (
	loggerContextKey loggerKeyType = iota
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

// ToContext set logger with context
func ToContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

// FromContext get logger from context
func FromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(loggerContextKey).(*zap.Logger); ok {
		return logger
	}
	return nil
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
