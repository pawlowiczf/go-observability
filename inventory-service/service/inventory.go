package service

import (
	"context"
	"time"

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

func (svc *Service) BurnCPU(ctx context.Context, duration time.Duration) (model.BurnCPUResponse, error) {
	deadline := time.Now().Add(duration)
	// Busy loop to consume CPU until the duration passes.
	var x float64
	for time.Now().Before(deadline) {
		x += 3.14159265 * 2.7182818
	}
	return model.BurnCPUResponse{Burned: true, Duration: duration.Seconds()}, nil
}
