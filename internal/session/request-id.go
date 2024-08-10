package session

import (
	"context"
)

// WithRequestID carries new context with rid request id.
func WithRequestID(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, requestIDKey, rid)
}

// RequestIDFromContext gets request id from context.
func RequestIDFromContext(ctx context.Context) (string, bool) {
	rid, ok := ctx.Value(requestIDKey).(string)
	return rid, ok
}
