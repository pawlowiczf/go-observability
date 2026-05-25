package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/pawlowiczf/go-observability/order-service/model"
	"github.com/pawlowiczf/go-observability/order-service/service"
	"github.com/pawlowiczf/go-observability/order-service/telemetry"
)

type Handler struct {
	svc   *service.Service
	ready atomic.Bool
}

func New(svc *service.Service) *Handler {
	h := &Handler{svc: svc}
	h.ready.Store(true)
	return h
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /checkout", h.Checkout)
	mux.HandleFunc("GET /healthz", h.Health)
	mux.HandleFunc("GET /readyz", h.Ready)
}

func (h *Handler) Checkout(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "request received",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
	)

	var req model.PlaceOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	start := time.Now()
	resp, err := h.svc.Checkout(r.Context(), req)
	success := err == nil
	duration := time.Since(start)

	if err != nil {
		slog.ErrorContext(r.Context(), "checkout failed", slog.Any("error", err))
		telemetry.RecordCheckoutRequest(r.Context(), duration, success)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	telemetry.RecordCheckoutRequest(r.Context(), duration, success)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Ready(w http.ResponseWriter, _ *http.Request) {
	if !h.ready.Load() {
		http.Error(w, "draining", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) NotReady() {
	h.ready.Store(false)
}
