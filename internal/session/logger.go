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
	logger := ctx.Value(loggingKey).(*slog.Logger)
	return logger
}
