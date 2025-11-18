package orders

import (
	"ecommerce-system/internal/json"
	"errors"
	"log"
	"net/http"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

type orderItem struct {
	ProductID int64 `json:"productId"`
	Quantity  int32 `json:"quantity"`
}

type createOrderRequest struct {
	CustomerID int64       `json:"customerId"`
	Items      []orderItem `json:"items"`
}

func (h *Handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var req createOrderRequest
	if err := json.Read(r, &req); err != nil {
		log.Printf("Error parsing request: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := h.service.PlaceOrder(r.Context(), req)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("Error placing order: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, order)
}
