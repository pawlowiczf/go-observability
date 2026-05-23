package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/pawlowiczf/go-observability/payment-service/model"
	"github.com/pawlowiczf/go-observability/payment-service/service"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /charge", h.Charge)
	mux.HandleFunc("GET /healthz", h.Health)
}

func (h *Handler) Charge(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "request received",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
	)

	var req model.ChargePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	resp, err := h.svc.Charge(r.Context(), req)
	if err != nil {
		slog.ErrorContext(r.Context(), "charge failed", slog.Any("error", err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
