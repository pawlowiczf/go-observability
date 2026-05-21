package config

import "github.com/caarlos0/env/v11"

type Config struct {
	ServiceName  string `env:"SERVICE_NAME" envDefault:"payment-service"`
	ServicePort  string `env:"SERVICE_PORT" envDefault:"8082"`
	OTELEndpoint string `env:"OTEL_ENDPOINT" envDefault:"localhost:4317"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	return cfg, err
}
