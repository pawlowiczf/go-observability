package client

import (
	"context"
	"fmt"
	"net/http"
)

type Payment struct {
	httpClient *http.Client
	baseURL    string
}

func NewPayment(httpClient *http.Client, baseURL string) *Payment {
	return &Payment{httpClient: httpClient, baseURL: baseURL}
}

type ChargePaymentRequest struct {
	Method    string  `json:"method"`
	CardLast4 string  `json:"card_last4"`
	Currency  string  `json:"currency"`
	Amount    float64 `json:"amount"`
}

type ChargePaymentResponse struct {
	Charged       bool   `json:"charged"`
	TransactionID string `json:"transaction_id"`
}

func (c *Payment) Charge(ctx context.Context, req ChargePaymentRequest) (ChargePaymentResponse, error) {
	var resp ChargePaymentResponse
	if err := postJSON(ctx, c.httpClient, c.baseURL+"/charge", req, &resp); err != nil {
		return ChargePaymentResponse{}, fmt.Errorf("payment charge: %w", err)
	}
	return resp, nil
}
