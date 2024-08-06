// HTTP server CLI endpoint for plans project. CLI - Command Line Interface.
package main

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"plans/cmd/primary-server/endpoint"
	"plans/cmd/primary-server/middleware"
	"time"
)

func main() {
	logHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(logHandler)
	router := mux.NewRouter()
	router.Use(middleware.Logging(logger))

	apiV1 := router.PathPrefix("/api/v1").Subrouter()
	endpoint.RegisterTest(apiV1)

	addr := "localhost:8080"

	server := http.Server{
		Addr:     addr,
		Handler:  router,
		ErrorLog: slog.NewLogLogger(logHandler, slog.LevelError),
	}

	ctxStop, cancelStop := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancelStop()

	g, gCtx := errgroup.WithContext(ctxStop)

	g.Go(func() error {
		logger.Info("Server is running on " + addr + "...")
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})
	g.Go(func() error {
		<-gCtx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(ctx)
	})

	if err := g.Wait(); err != nil {
		logger.Error("Failed serving: " + err.Error())
	}

	logger.Info("Server gracefully shutdown")
}
