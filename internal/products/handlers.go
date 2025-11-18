package products

import (
	"ecommerce-system/internal/json"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, products)
}

func (h *Handler) FindProductByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "product ID must be require", http.StatusBadRequest)
		return
	}
	productID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "product ID must be a number", http.StatusBadRequest)
	}

	product, err := h.service.FindProductByID(r.Context(), productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.Write(w, http.StatusOK, product)
}
