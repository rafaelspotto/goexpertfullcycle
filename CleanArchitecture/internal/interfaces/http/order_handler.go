package http

import (
	"encoding/json"
	"net/http"

	"github.com/rafaelspotto/goexpertfullcycle/cleanarchitecture/internal/core/usecase"
)

type OrderHandler struct {
	OrderUseCase *usecase.OrderUseCase
}

func NewOrderHandler(useCase *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		OrderUseCase: useCase,
	}
}

type CreateOrderRequest struct {
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request CreateOrderRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	order, err := h.OrderUseCase.Create(request.Price, request.Tax)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	orders, err := h.OrderUseCase.List()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
