package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type requestID struct{}

// WithRequestID injects request id into the context.
// Generates new UUID if id is empty string.
// Generates new UUID if id isn't a correct UUID.
func WithRequestID(ctx context.Context, id string) context.Context {
	err := uuid.Validate(id)
	if id == "" || err != nil {
		id = uuid.NewString()
	}

	return context.WithValue(ctx, requestID{}, id)
}

// RequestIDFrom extracts request id from context.
// Returns empty string if request id can not be found.
func RequestIDFrom(ctx context.Context) string {
	reqID, _ := ctx.Value(requestID{}).(string)
	return reqID
}

// RequestID middleware injects request id for each request.
// It searches for X-Request-ID header in request, if not present, generates a new one.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		r = r.WithContext(WithRequestID(r.Context(), reqID))
		next.ServeHTTP(w, r)
	})
}
