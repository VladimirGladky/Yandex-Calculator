package logger

import (
	"context"
	"go.uber.org/zap"
)

const (
	Key = "logger"
)

type Logger struct {
	l *zap.Logger
}

func (l Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.l.Info(msg, fields...)
}

func (l Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.l.Error(msg, fields...)
}

func New(ctx context.Context) (context.Context, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, Key, &Logger{l: logger})
	return ctx, nil
}

func GetLoggerFromCtx(ctx context.Context) *Logger {
	return ctx.Value(Key).(*Logger)
}
