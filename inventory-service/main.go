package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pawlowiczf/go-observability/inventory-service/config"
	"github.com/pawlowiczf/go-observability/inventory-service/handler"
	"github.com/pawlowiczf/go-observability/inventory-service/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("could not load config", slog.Any("error", err))
		os.Exit(1)
	}

	svc := service.New()

	mux := http.NewServeMux()

	h := handler.New(svc)
	h.Register(mux)

	srv := &http.Server{
		Addr:    ":" + cfg.ServicePort,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", slog.Any("error", err))
			os.Exit(1)
		}
	}()
	slog.Info("server started", slog.String("addr", srv.Addr))

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown error", slog.Any("error", err))
	}
}
