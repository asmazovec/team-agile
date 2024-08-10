package session

import (
	"context"
	"log/slog"
)

// WithLogger carries new context with configured slog.Logger value.
func WithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, loggingKey, l)
}

// MustLoggerFromContext gets session slog.Logger from context.
func MustLoggerFromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggingKey).(*slog.Logger)
	if !ok {
		panic("session: logger not found in context")
	}
	return logger
}
