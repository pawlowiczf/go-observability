
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

Kolejność implementacji (od liścia w górę)
1. inventory-service i payment-service (mogą iść równolegle, identyczny pattern):

model/inventory.go — ReserveRequest{ProductID, Quantity}, ReserveResponse{Reserved bool}.
service/inventory.go — dorzucam metodę Reserve(ctx, req) (*ReserveResponse, error). Na razie zwraca {Reserved: true}.
handler/inventory.go — Handler struct z *service.Service, metoda Reserve(w, r) robi json.Decode → service.Reserve → json.Encode. Plus Health(w, r).
main.go — mux.HandleFunc("POST /reserve", ...), mux.HandleFunc("GET /healthz", ...), http.Server, graceful shutdown.
2. order-service (zależy od 1, bo używa kontraktów):

model/checkout.go — CheckoutRequest z requests.md (bez pola simulate na razie).
client/inventory.go — własne ReserveRequest/Response (kopia), Client z *http.Client + baseURL z config, metoda Reserve(ctx, ...).
client/payment.go — analogicznie.
service/order.go — Service dostaje obu clientów, ma Checkout(ctx, req) które woła inventory.Reserve → payment.Charge → zwraca status.
handler/order.go — POST /checkout.
main.go — j.w. plus konstrukcja klientów z URL-ami z config.