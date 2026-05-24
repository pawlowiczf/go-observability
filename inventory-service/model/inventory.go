package model

type ReserveProductRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type ReserveProductResponse struct {
	Reserved bool `json:"reserved"`
}

type BurnCPUResponse struct {
	Burned   bool    `json:"burned"`
	Duration float64 `json:"duration_seconds"`
}
