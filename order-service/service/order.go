package service

import (
	"context"

	"github.com/pawlowiczf/go-observability/order-service/model"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (svc *Service) Checkout(ctx context.Context, req model.PlaceOrderRequest) (model.PlaceOrderResponse, error) {
	return model.PlaceOrderResponse{
		OrderID: "order-stub-0001",
		Status:  "confirmed",
	}, nil
}
