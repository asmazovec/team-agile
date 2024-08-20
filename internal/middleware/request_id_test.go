package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asmazovec/team-agile/internal/middleware"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestWithRequestID_ShouldInjectRequestID(t *testing.T) {
	id := uuid.NewString()
	ctx := context.Background()

	ctx = middleware.WithRequestID(ctx, id)
	res := middleware.RequestIDFrom(ctx)

	assert.Equal(t, id, res)
}

func TestRequestIDFrom_WithEmptyContext_ShouldEmptyString(t *testing.T) {
	ctx := context.Background()

	res := middleware.RequestIDFrom(ctx)

	assert.Equal(t, "", res)
}

func TestWithRequestID_EmptyID_ShouldGenerateNewID(t *testing.T) {
	id := ""
	ctx := context.Background()

	ctx = middleware.WithRequestID(ctx, id)
	res := middleware.RequestIDFrom(ctx)

	assert.NoError(t, uuid.Validate(res))
}

func TestWithRequestID_WrongID_ShouldGenerateNewID(t *testing.T) {
	id := "123asdf"
	ctx := context.Background()

	ctx = middleware.WithRequestID(ctx, id)
	res := middleware.RequestIDFrom(ctx)

	assert.NoError(t, uuid.Validate(res))
}

func TestRequestID_ServeHTTP_ShouldAddRequestID(t *testing.T) {
	nextHandler := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		reqID := middleware.RequestIDFrom(r.Context())

		assert.NotEqual(t, "", reqID)
		assert.NoError(t, uuid.Validate(reqID))
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	mw := middleware.RequestID(nextHandler)

	mw.ServeHTTP(httptest.NewRecorder(), req)
}

func TestRequestID_ServeHTTP_WithRequestIDInRequest_ShouldNotGenerate(t *testing.T) {
	id := uuid.NewString()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-ID", id)
	nextHandler := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		reqID := middleware.RequestIDFrom(r.Context())

		assert.Equal(t, id, reqID)
	})
	mw := middleware.RequestID(nextHandler)

	mw.ServeHTTP(httptest.NewRecorder(), req)
}
