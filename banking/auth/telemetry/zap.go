package telemetry

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

type loggerKey struct{}

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func Logger(ctx context.Context) (*zap.Logger, error) {
	if ctx == nil {
		return nil, errors.New("nil context")
	}

	logger, ok := ctx.Value(loggerKey{}).(*zap.Logger)
	if !ok {
		return nil, errors.New("logger not found")
	}

	return logger, nil
}

func NewNopLogger() *zap.Logger {
	return zap.NewNop()
}
