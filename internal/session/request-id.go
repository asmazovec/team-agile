package session

import (
	"context"
)

// WithRequestId carries new context with rid request id.
func WithRequestId(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, requestIdKey, rid)
}

// RequestIdFromContext gets request id from context.
func RequestIdFromContext(ctx context.Context) (string, bool) {
	rid, ok := ctx.Value(requestIdKey).(string)
	return rid, ok
}
