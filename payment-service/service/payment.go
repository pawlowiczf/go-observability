package service

import (
	"context"
	"errors"
	"math/rand"

	"github.com/pawlowiczf/go-observability/payment-service/model"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (svc *Service) Charge(ctx context.Context, req model.ChargePaymentRequest) (model.ChargePaymentResponse, error) {
	// Simulate occasional payment failures so we can observe error traces.
	if rand.Intn(5) == 0 {
		return model.ChargePaymentResponse{}, errors.New("payment declined")
	}

	return model.ChargePaymentResponse{
		Charged:       true,
		TransactionID: "txn-stub-0001",
	}, nil
}
