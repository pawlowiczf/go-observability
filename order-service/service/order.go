package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pawlowiczf/go-observability/order-service/client"
	"github.com/pawlowiczf/go-observability/order-service/model"
)

type Service struct {
	inventory *client.Inventory
	payment   *client.Payment
}

func New(inventory *client.Inventory, payment *client.Payment) *Service {
	return &Service{
		inventory: inventory,
		payment:   payment,
	}
}

func (s *Service) Checkout(ctx context.Context, req model.PlaceOrderRequest) (model.PlaceOrderResponse, error) {
	var total float64
	for _, it := range req.Items {
		resp, err := s.inventory.Reserve(ctx, client.ReserveProductRequest{
			ProductID: it.ProductID,
			Quantity:  it.Quantity,
		})
		if err != nil {
			return model.PlaceOrderResponse{}, fmt.Errorf("reserve %s: %w", it.ProductID, err)
		}
		if !resp.Reserved {
			return model.PlaceOrderResponse{}, fmt.Errorf("reserve %s: not reserved", it.ProductID)
		}
		total += float64(it.Quantity) * it.UnitPrice
	}

	charge, err := s.payment.Charge(ctx, client.ChargePaymentRequest{
		Method:    req.Payment.Method,
		CardLast4: req.Payment.CardLast4,
		Currency:  req.Payment.Currency,
		Amount:    total,
	})
	if err != nil {
		return model.PlaceOrderResponse{}, fmt.Errorf("charge: %w", err)
	}
	if !charge.Charged {
		return model.PlaceOrderResponse{}, fmt.Errorf("charge declined")
	}

	return model.PlaceOrderResponse{
		OrderID: "order-" + uuid.NewString(),
		Status:  "confirmed",
	}, nil
}
