package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/asmazovec/team-agile/internal/closer"
	"github.com/asmazovec/team-agile/internal/config"
	mw "github.com/asmazovec/team-agile/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// main application lifecycle entry point.
func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "c", "", "Path to configuration file")
	flag.Parse()

	cfg := config.MustRead(config.FromEnv(cfgPath))
	l := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	c := &closer.Closer{}

	err := run(c, l, cfg)
	if err != nil {
		l.Error(err.Error())
		panic(err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.AppShutdownTimeout)
	defer shutdownCancel()

	l.Info("Shutting down")
	errCh := c.Close(shutdownCtx)
	for err := range errCh {
		if err != nil {
			l.Error(fmt.Sprintf("Shutting down: %v", err.Error()))
		}
	}
}

//nolint:unparam // currently not needed
func run(c *closer.Closer, l *slog.Logger, cfg config.AppConfig) error {
	srv := &http.Server{
		Addr:              cfg.HTTPPrimaryServer.Address,
		ReadTimeout:       cfg.HTTPPrimaryServer.ReadTimeout,
		ReadHeaderTimeout: cfg.HTTPPrimaryServer.ReadHeaderTimeout,
	}

	router := chi.NewRouter()
	router.Use(render.SetContentType(render.ContentTypeJSON))
	router.Use(mw.RequestID)
	router.Use(mw.Logger(l, mw.RequestIDLog, mw.MethodLog))

	router.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello world!"))
	})

	srv.Handler = router

	// Starting resources ...
	go func() {
		l.Info("Starting server on " + cfg.HTTPPrimaryServer.Address)
		_ = srv.ListenAndServe()
	}()
	_, _ = c.Add(closer.ReleaserWithLog(l, "Closing HTTP Primary server", srv.Shutdown))

	return nil
}
