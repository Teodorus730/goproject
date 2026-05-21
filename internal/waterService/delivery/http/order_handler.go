package handler

import (
    "encoding/json"
    "net/http"
    "log/slog"
	"context"
	"time"

    "goproject/internal/domain"

    "github.com/go-playground/validator/v10"
)

type OrderHandler struct {
	uc domain.OrderUsecase
	log *slog.Logger
}

func NewOrderHandler(uc domain.OrderUsecase, log *slog.Logger) *OrderHandler {
	return &OrderHandler{uc: uc, log: log}
}

var validate = validator.New()
const CtxTimeout = 5 * time.Second


func (h *OrderHandler) AddNewOrder(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Received new order request")
	h.log.Info("Incoming auth header", slog.String("header", r.Header.Get("Authorization")))

	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		h.log.Error("User not authenticated, missing userID")
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var req struct {
		TariffID int    `json:"tarif_id"`
		Amount   int    `json:"amount"`
		Address  string `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode request body", slog.Any("error", err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	order, err := h.uc.PlaceOrder(r.Context(), userID, req.TariffID, req.Amount, req.Address)
	if err != nil {
		h.log.Error("PlaceOrder failed", slog.Any("error", err))
		http.Error(w, "failed to place order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(order); err != nil {
		h.log.Error("Failed to encode response", slog.Any("error", err))
	}
}

func (h *OrderHandler) GetOrdersList(w http.ResponseWriter, r *http.Request) {
    h.log.Info("Received get orders list request")
	h.log.Info("Incoming auth header", slog.String("header", r.Header.Get("Authorization")))

	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		h.log.Error("User not authenticated, missing userID")
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

    _, getOrdersCancel := context.WithTimeout(context.Background(), CtxTimeout)
    defer getOrdersCancel()
    orders, err := h.uc.GetOrders(r.Context(), userID)
    if err != nil {
		h.log.Error("Failed to get orders", slog.Any("error", err))
        return
    }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		h.log.Error("Failed to encode response", slog.Any("error", err))
	}
}