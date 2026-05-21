package config

import "github.com/caarlos0/env/v11"

type Config struct {
	ServiceName      string `env:"SERVICE_NAME" envDefault:"order-service"`
	ServicePort      string `env:"SERVICE_PORT" envDefault:"8081"`
	OTELEndpoint     string `env:"OTEL_ENDPOINT" envDefault:"localhost:4317"`
	InventoryBaseURL string `env:"INVENTORY_BASE_URL" envDefault:"http://localhost:8080"`
	PaymentBaseURL   string `env:"PAYMENT_BASE_URL" envDefault:"http://localhost:8082"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	return cfg, err
}
