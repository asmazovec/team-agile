package middleware

import (
	"context"
	"log/slog"
	"net/http"
)

type loggerKey struct{}

// WithLogger injects logger in the context.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	ctx = context.WithValue(ctx, loggerKey{}, logger)
	return ctx
}

// LoggerFrom gets a logger from the context.
func LoggerFrom(ctx context.Context) *slog.Logger {
	l, _ := ctx.Value(loggerKey{}).(*slog.Logger)
	return l
}

type logger struct {
	h http.Handler
	l *slog.Logger
	e []LogExtension
}

// LogExtension is an additional attributes for logging.
type LogExtension func(r *http.Request) (slog.Attr, bool)

// Logger middleware injects logger into each request context
// Logger is configured slog.Logger with additions of LogExtension.
func Logger(l *slog.Logger, extensions ...LogExtension) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return &logger{next, l, extensions}
	}
}

// RequestIDLog logger extension for log request id.
func RequestIDLog(r *http.Request) (slog.Attr, bool) {
	reqID := RequestIDFrom(r.Context())
	if reqID == "" {
		return slog.Attr{}, false
	}
	return slog.String("request-id", reqID), true
}

// MethodLog logger extension for log request method.
func MethodLog(r *http.Request) (slog.Attr, bool) {
	return slog.String("method", r.Method), true
}

// ServeHTTP implements http.Handler interface.
func (l *logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lg := l.l
	for _, ext := range l.e {
		if ext == nil {
			continue
		}

		a, ok := ext(r)
		if !ok {
			continue
		}

		lg = lg.With(a)
	}

	lg.Info("Request")

	ctx := WithLogger(r.Context(), lg)
	r = r.WithContext(ctx)
	l.h.ServeHTTP(w, r)
}
