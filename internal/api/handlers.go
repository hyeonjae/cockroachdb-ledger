package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mini-ledger/internal/domain"
	"mini-ledger/internal/service"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	tradingService *service.TradingService
}

func NewHandler(tradingService *service.TradingService) *Handler {
	return &Handler{
		tradingService: tradingService,
	}
}

func (h *Handler) GetAccountBalance(w http.ResponseWriter, r *http.Request) {
	accountIDStr := chi.URLParam(r, "accountID")
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		h.writeErrorResponse(w, "invalid account ID", http.StatusBadRequest)
		return
	}

	balance, err := h.tradingService.GetAccountBalance(accountID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.writeJSONResponse(w, balance, http.StatusOK)
}

func (h *Handler) GetAccountHoldings(w http.ResponseWriter, r *http.Request) {
	accountIDStr := chi.URLParam(r, "accountID")
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		h.writeErrorResponse(w, "invalid account ID", http.StatusBadRequest)
		return
	}

	holdings, err := h.tradingService.GetAccountHoldings(accountID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.writeJSONResponse(w, holdings, http.StatusOK)
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "invalid request body", http.StatusBadRequest)
		return
	}

	order, err := h.tradingService.CreateOrder(&req)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.writeJSONResponse(w, order, http.StatusCreated)
}

func (h *Handler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	orderIDStr := chi.URLParam(r, "orderID")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		h.writeErrorResponse(w, "invalid order ID", http.StatusBadRequest)
		return
	}

	order, err := h.tradingService.CancelOrder(orderID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.writeJSONResponse(w, order, http.StatusOK)
}

func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrAccountNotFound:
		h.writeErrorResponse(w, "account not found", http.StatusNotFound)
	case domain.ErrOrderNotFound:
		h.writeErrorResponse(w, "order not found", http.StatusNotFound)
	case domain.ErrInsufficientFunds:
		h.writeErrorResponse(w, "insufficient funds", http.StatusBadRequest)
	case domain.ErrInsufficientHoldingQuantity:
		h.writeErrorResponse(w, "insufficient holding quantity", http.StatusBadRequest)
	case domain.ErrOrderNotCancelable:
		h.writeErrorResponse(w, "order is not in a cancelable state", http.StatusBadRequest)
	default:
		h.writeErrorResponse(w, "internal server error", http.StatusInternalServerError)
	}
}

func (h *Handler) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(domain.ErrorResponse{Error: message})
}