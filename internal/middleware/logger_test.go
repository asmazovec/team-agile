package middleware_test

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/asmazovec/team-agile/internal/middleware"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestWithLogger_ShouldLogger(t *testing.T) {
	l := &slog.Logger{}
	ctx := context.Background()

	ctx = middleware.WithLogger(ctx, l)
	res := middleware.LoggerFrom(ctx)

	assert.Equal(t, l, res)
}

func TestLoggerFrom_WithEmptyCtx_ShouldNil(t *testing.T) {
	ctx := context.Background()

	res := middleware.LoggerFrom(ctx)

	assert.Nil(t, res)
}

func TestLogger_ServeHTTP_ShouldAddLogger(t *testing.T) {
	f, _ := os.Open(os.DevNull)
	defer f.Close()
	origLogger := slog.New(slog.NewJSONHandler(f, nil))
	nextHandler := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		l := middleware.LoggerFrom(r.Context())
		assert.Equal(t, origLogger, l)
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	mw := middleware.Logger(origLogger)(nextHandler)

	mw.ServeHTTP(httptest.NewRecorder(), req)
}

func TestLogger_ServeHTTP_NilExt_ShouldNotPanic(t *testing.T) {
	f, _ := os.Open(os.DevNull)
	defer f.Close()
	origLogger := slog.New(slog.NewJSONHandler(f, nil))
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	mw := middleware.Logger(origLogger, nil)(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))

	assert.NotPanics(t, func() {
		mw.ServeHTTP(httptest.NewRecorder(), req)
	})
}

func TestMethodLog_ShouldLogMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodConnect, "/", nil)

	arg, ok := middleware.MethodLog(req)

	assert.True(t, ok)
	assert.Equal(t, req.Method, arg.Value.String())
	assert.Equal(t, "method", arg.Key)
}

func TestRequestIDLog_ShouldLogRequestID(t *testing.T) {
	id := uuid.NewString()
	req := httptest.NewRequest(http.MethodConnect, "/", nil)
	req = req.WithContext(middleware.WithRequestID(req.Context(), id))

	arg, ok := middleware.RequestIDLog(req)

	assert.True(t, ok)
	assert.Equal(t, id, arg.Value.String())
	assert.Equal(t, "request-id", arg.Key)
}

func TestRequestIDLog_WithNoRequestIDInContext_ShouldNotOk(t *testing.T) {
	req := httptest.NewRequest(http.MethodConnect, "/", nil)

	_, ok := middleware.RequestIDLog(req)

	assert.False(t, ok)
}

type LogExtMock struct {
	called bool
}

func (lem *LogExtMock) ShouldCalled(_ *http.Request) (slog.Attr, bool) {
	lem.called = true
	return slog.Attr{}, true
}

func TestLogger_ServeHTTP_Ext_ShouldCall(t *testing.T) {
	f, _ := os.Open(os.DevNull)
	defer f.Close()
	origLogger := slog.New(slog.NewJSONHandler(f, nil))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ext := new(LogExtMock)
	mw := middleware.Logger(origLogger, ext.ShouldCalled)(
		http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))

	mw.ServeHTTP(httptest.NewRecorder(), req)

	assert.True(t, ext.called)
}
