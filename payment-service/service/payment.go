package service

import "github.com/pawlowiczf/go-observability/payment-service/config"

type Service struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Service {
	return &Service{
		cfg: cfg,
	}
}
