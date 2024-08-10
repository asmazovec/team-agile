package middleware

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/asmazovec/team-agile/internal/session"
)

type logging struct {
	h http.Handler
	l *slog.Logger
}

// Logging provides Logging middleware which includes logger to request context.
func Logging(l *slog.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &logging{l: l, h: h}
	}
}

func (l *logging) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rid := r.Header.Get("X-Request-ID")
	if rid == "" {
		rid = uuid.NewString()
	}

	urlLog := slog.String("url", r.URL.String())
	methodLog := slog.String("method", r.Method)
	ridLog := slog.String("request-id", rid)
	l.l = l.l.With(urlLog, ridLog, methodLog)

	r = r.WithContext(session.WithRequestID(r.Context(), rid))
	r = r.WithContext(session.WithLogger(r.Context(), l.l))

	l.l.Info("Request")
	l.h.ServeHTTP(w, r)
}
