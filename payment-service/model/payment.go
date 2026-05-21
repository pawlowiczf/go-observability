package model

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
