package service

import (
	"context"

	"github.com/pawlowiczf/go-observability/inventory-service/model"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (svc *Service) Reserve(ctx context.Context, req model.ReserveProductRequest) (model.ReserveProductResponse, error) {
	return model.ReserveProductResponse{Reserved: true}, nil
}
