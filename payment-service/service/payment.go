package service

import (
	"context"

	"github.com/pawlowiczf/go-observability/payment-service/model"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (svc *Service) Charge(ctx context.Context, req model.ChargePaymentRequest) (model.ChargePaymentResponse, error) {
	return model.ChargePaymentResponse{
		Charged:       true,
		TransactionID: "txn-stub-0001",
	}, nil
}
