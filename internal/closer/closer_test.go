package closer_test

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/asmazovec/team-agile/internal/closer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdd_Single_ShouldRegister(t *testing.T) {
	c := &closer.Closer{}

	res, err := c.Add(nil)

	assert.NotNilf(t, res, "Should register resource.")
	assert.NoError(t, err)
}

func TestAdd_Dependency_ShouldRegister(t *testing.T) {
	c := &closer.Closer{}
	r, _ := c.Add(nil)

	res, err := c.Add(nil, r)

	assert.NotNilf(t, res, "Should register resource")
	assert.NoError(t, err)
}

func TestAdd_MultipleDeps_ShouldRegisterAll(t *testing.T) {
	c := &closer.Closer{}
	r1, _ := c.Add(nil)
	r2, _ := c.Add(nil, r1)
	res, err := c.Add(nil, r1, r2)

	assert.NotNilf(t, res, "Should register resource")
	assert.NoError(t, err)
}

func TestAdd_SameMultipleTimes_ShouldRegisterOnce(t *testing.T) {
	c := &closer.Closer{}
	r1, _ := c.Add(nil)

	res, err := c.Add(nil, r1, r1, r1)

	assert.NotNilf(t, res, "Should register resource")
	assert.NoError(t, err)
}

func TestAdd_NilDependency_ShouldError(t *testing.T) {
	c := &closer.Closer{}

	res, err := c.Add(nil, nil, nil)

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestAdd_NotAssociatedDeps_ShouldError(t *testing.T) {
	c1 := &closer.Closer{}
	c2 := &closer.Closer{}
	r, _ := c1.Add(nil)

	res, err := c2.Add(nil, r)
	assert.Nil(t, res)
	assert.Error(t, err)
}

type ResourceMock struct {
	mu    sync.Mutex
	Order []int
}

func (r *ResourceMock) CallOrdered(order int, err error) closer.Releaser {
	return func(_ context.Context) error {
		r.mu.Lock()
		r.Order = append(r.Order, order)
		r.mu.Unlock()
		return err
	}
}

func (r *ResourceMock) CallOrderedWithTimeout(timeout time.Duration, order int, err error) closer.Releaser {
	ctx, closeCtx := context.WithTimeout(context.Background(), timeout)
	return func(baseCtx context.Context) error {
		defer closeCtx()
		select {
		case <-ctx.Done():
			r.mu.Lock()
			r.Order = append(r.Order, order)
			r.mu.Unlock()
		case <-baseCtx.Done():
		}
		return err
	}
}

func TestCancel_ShouldAwaitReleases(t *testing.T) {
	var (
		r = new(ResourceMock)
		c = new(closer.Closer)
	)

	r1, _ := c.Add(r.CallOrderedWithTimeout(100*time.Millisecond, 1, nil))
	r2, _ := c.Add(r.CallOrderedWithTimeout(10*time.Millisecond, 2, nil), r1)
	_, _ = c.Add(r.CallOrderedWithTimeout(80*time.Millisecond, 3, nil), r2, r1)
	_, _ = c.Add(r.CallOrderedWithTimeout(30*time.Millisecond, 4, nil), r2)
	errs := c.Close(context.Background())
	for err := range errs {
		require.NoError(t, err)
	}

	assert.True(t, slices.Equal(r.Order, []int{4, 3, 2, 1}))
}

func TestCancel_NilReleaser_ShouldNotPanic(_ *testing.T) {
	var c = new(closer.Closer)

	_, _ = c.Add(nil)
	errs := c.Close(context.Background())

	<-errs
}

func TestCancel_ShouldReleaseEverySubgraph(t *testing.T) {
	var (
		r = new(ResourceMock)
		c = new(closer.Closer)
	)

	// Subgraph 1
	g1r1, _ := c.Add(r.CallOrdered(3, nil))
	g1r2, _ := c.Add(r.CallOrdered(2, nil), g1r1)
	_, _ = c.Add(r.CallOrdered(1, nil), g1r2, g1r1)

	// Subgraph 2
	g2r1, _ := c.Add(r.CallOrdered(2, nil))
	_, _ = c.Add(r.CallOrdered(1, nil), g2r1)
	_, _ = c.Add(r.CallOrdered(1, nil), g2r1)

	// Subgraph 3
	_, _ = c.Add(r.CallOrdered(1, nil))

	errs := c.Close(context.Background())
	for err := range errs {
		require.NoError(t, err)
	}

	assert.True(t, slices.Equal(r.Order, []int{1, 1, 1, 1, 2, 2, 3}))
}

func TestCancel_LayerShouldReleaseInParallel(t *testing.T) {
	var (
		r = new(ResourceMock)
		c = new(closer.Closer)
	)
	ctx, done := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer done()

	r1, _ := c.Add(r.CallOrdered(1, nil))
	_, _ = c.Add(r.CallOrderedWithTimeout(100*time.Millisecond, 2, nil), r1)
	_, _ = c.Add(r.CallOrderedWithTimeout(100*time.Millisecond, 2, nil), r1)
	_, _ = c.Add(r.CallOrderedWithTimeout(100*time.Millisecond, 2, nil), r1)
	errs := c.Close(ctx)
	for err := range errs {
		require.NoError(t, err)
	}

	assert.True(t, slices.Equal(r.Order, []int{2, 2, 2, 1}))
}

func TestCancel_ExpireContext_ShouldStop(t *testing.T) {
	var (
		r = new(ResourceMock)
		c = new(closer.Closer)
	)
	ctx, cancel := context.WithCancel(context.Background())
	r1, _ := c.Add(r.CallOrderedWithTimeout(10*time.Millisecond, 1, nil))
	_, _ = c.Add(r.CallOrdered(2, nil), r1)

	errs := c.Close(ctx)
	<-errs
	cancel()
	<-errs

	assert.True(t, slices.Equal(r.Order, []int{2}))
}

func TestCancel_Graph_ShouldCancelInOrder(t *testing.T) {
	var (
		r = new(ResourceMock)
		c = new(closer.Closer)
	)

	r1, _ := c.Add(r.CallOrdered(1, nil))
	r2, _ := c.Add(r.CallOrdered(2, nil), r1)
	_, _ = c.Add(r.CallOrdered(3, nil), r2, r1)
	_, _ = c.Add(r.CallOrdered(3, nil), r2)
	errs := c.Close(context.Background())
	for err := range errs {
		require.NoError(t, err)
	}

	assert.True(t, slices.Equal(r.Order, []int{3, 3, 2, 1}))
}

func TestCancel_GraphWithErrorCall_ShouldAddToErrorsMsg(t *testing.T) {
	var (
		r = new(ResourceMock)
		c = new(closer.Closer)
	)

	errExpected := errors.New("error")
	r1, _ := c.Add(r.CallOrdered(1, errExpected))
	r2, _ := c.Add(r.CallOrdered(2, errExpected), r1)
	_, _ = c.Add(r.CallOrdered(3, nil), r2)
	_, _ = c.Add(r.CallOrdered(3, errExpected), r2)
	errs := c.Close(context.Background())

	errCount := 0
	for err := range errs {
		if err != nil {
			errCount++
			require.ErrorIs(t, err, errExpected)
		}
	}
	assert.Equal(t, 3, errCount)
}

type LogHandlerMock struct {
	called bool
}

func (lhm *LogHandlerMock) Enabled(context.Context, slog.Level) bool  { return true }
func (lhm *LogHandlerMock) Handle(context.Context, slog.Record) error { lhm.called = true; return nil }
func (lhm *LogHandlerMock) WithAttrs([]slog.Attr) slog.Handler        { return lhm }
func (lhm *LogHandlerMock) WithGroup(string) slog.Handler             { return lhm }

func TestReleaserWithLog_ShouldCallLog(t *testing.T) {
	h := &LogHandlerMock{}
	l := slog.New(h)
	r := closer.ReleaserWithLog(l, "", func(context.Context) error { return nil })

	_ = r(nil)

	assert.True(t, h.called)
}

func TestReleaserWithLog_CallNilReleaser_ShouldNotPanic(t *testing.T) {
	h := &LogHandlerMock{}
	l := slog.New(h)
	r := closer.ReleaserWithLog(l, "", nil)
	assert.NotPanics(t, func() {
		_ = r(nil)
	})
}

func TestReleaserWithLog_ShouldCallOriginalReleaser(t *testing.T) {
	h := &LogHandlerMock{}
	l := slog.New(h)
	res := &ResourceMock{}
	r := closer.ReleaserWithLog(l, "", res.CallOrdered(0, nil))

	_ = r(nil)

	assert.Len(t, res.Order, 1)
}
