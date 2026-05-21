package main

import (
	"log/slog"
	"os"

	"github.com/pawlowiczf/go-observability/order-service/config"
	"github.com/pawlowiczf/go-observability/order-service/service"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error(
			"could not load config",
			slog.Any("error", err),
		)
	}

	service.New(cfg)
}
