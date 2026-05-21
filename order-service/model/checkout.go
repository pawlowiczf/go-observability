package model

type Item struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

type Payment struct {
	Method    string `json:"method"`
	CardLast4 string `json:"card_last4"`
	Currency  string `json:"currency"`
}

type PlaceOrderRequest struct {
	CustomerID string  `json:"customer_id"`
	SessionID  string  `json:"session_id"`
	Items      []Item  `json:"items"`
	Payment    Payment `json:"payment"`
}

type PlaceOrderResponse struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}
