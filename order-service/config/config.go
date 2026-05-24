package config

import "github.com/caarlos0/env/v11"

type Config struct {
	ServiceName      string `env:"SERVICE_NAME" envDefault:"order-service"`
	ServicePort      string `env:"SERVICE_PORT" envDefault:"8080"`
	OTELEndpoint     string `env:"OTEL_ENDPOINT" envDefault:"otel-collector:4318"`
	InventoryBaseURL string `env:"INVENTORY_BASE_URL" envDefault:"inventory-service:8080"`
	PaymentBaseURL   string `env:"PAYMENT_BASE_URL" envDefault:"payment-service:8080"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	return cfg, err
}
