package client

import (
	"context"
	"fmt"
	"net/http"
)

type Inventory struct {
	httpClient *http.Client
	baseURL    string
}

func NewInventory(httpClient *http.Client, baseURL string) *Inventory {
	return &Inventory{httpClient: httpClient, baseURL: baseURL}
}

type ReserveProductRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type ReserveProductResponse struct {
	Reserved bool `json:"reserved"`
}

func (c *Inventory) Reserve(ctx context.Context, req ReserveProductRequest) (ReserveProductResponse, error) {
	var resp ReserveProductResponse
	if err := postJSON(ctx, c.httpClient, c.baseURL+"/reserve", req, &resp); err != nil {
		return ReserveProductResponse{}, fmt.Errorf("inventory reserve: %w", err)
	}
	return resp, nil
}
