package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/pawlowiczf/go-observability/inventory-service/model"
	"github.com/pawlowiczf/go-observability/inventory-service/service"
	"github.com/pawlowiczf/go-observability/inventory-service/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /reserve", h.Reserve)
	mux.HandleFunc("GET /healthz", h.Health)
}

func (h *Handler) Reserve(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "request received",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
	)
	span := trace.SpanFromContext(r.Context())

	var req model.ReserveProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	span.SetAttributes(
		attribute.String("product_id", req.ProductID),
		attribute.Int("quantity", req.Quantity),
	)

	start := time.Now()
	resp, err := h.svc.Reserve(r.Context(), req)
	success := err == nil
	telemetry.RecordReserveRequest(r.Context(), time.Since(start), success)

	if err != nil {
		slog.ErrorContext(r.Context(), "reserve failed", slog.Any("error", err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
