package service

import "github.com/pawlowiczf/go-observability/order-service/config"

type Service struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Service {
	return &Service{
		cfg: cfg,
	}
}
