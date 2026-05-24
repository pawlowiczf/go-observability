package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pawlowiczf/go-observability/inventory-service/config"
	"github.com/pawlowiczf/go-observability/inventory-service/handler"
	"github.com/pawlowiczf/go-observability/inventory-service/service"
	"github.com/pawlowiczf/go-observability/inventory-service/telemetry"
	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/sync/errgroup"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("could not load config", slog.Any("error", err))
		return 1
	}

	shutdown, err := telemetry.Setup(ctx, cfg.ServiceName, cfg.OTELEndpoint)
	if err != nil {
		slog.Error("could not setup telemetry", slog.Any("error", err))
		return 1
	}
	if err := telemetry.InitMetrics(cfg.ServiceName); err != nil {
		slog.Error("could not init metrics", slog.Any("error", err))
		return 1
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdown(shutdownCtx); err != nil {
			slog.Error("telemetry shutdown error", slog.Any("error", err))
		}
	}()

	slog.SetDefault(slog.New(slogmulti.Fanout(
		slog.NewJSONHandler(os.Stdout, nil),
		otelslog.NewHandler(cfg.ServiceName),
	)))

	svc := service.New()
	mux := http.NewServeMux()
	h := handler.New(svc)
	h.Register(mux)
	srv := &http.Server{
		Addr:    ":" + cfg.ServicePort,
		Handler: otelhttp.NewHandler(mux, cfg.ServiceName),
	}

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		slog.Info("server started", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http server: %w", err)
		}
		return nil
	})
	g.Go(func() error {
		<-gCtx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	})

	if err := g.Wait(); err != nil {
		slog.Error("shutdown with error", slog.Any("error", err))
		return 1
	}
	return 0
}
