
```sh
POST /checkout
{
  "customer_id": "cust-42",
  "session_id":  "sess-abc123",
  "items": [
    {
      "product_id": "prod-001",
      "name":       "Mechanical Keyboard",
      "quantity":   2,
      "unit_price": 299.99
    }
  ],
  "payment": {
    "method":    "credit_card",
    "card_last4": "4242",
    "currency":  "PLN"
  },
  "simulate": {
    "inventory_out_of_stock": false,
    "inventory_delay_ms":     0,
    "payment_fail":           false,
    "payment_delay_ms":       0,
    "order_fail":             false
  }
}
```