package handler

import (
	"encoding/json"
	"net/http"

	"github.com/farhaan/kuncie/entity"
	"github.com/farhaan/kuncie/internal/checkout"
)

type checkoutHandler struct {
	checkoutSvc checkout.Service
}
type CheckoutHandler interface {
	Checkout(rw http.ResponseWriter, req *http.Request)
}

func NewCheckoutHandler(checkoutSvc checkout.Service) CheckoutHandler {
	return &checkoutHandler{checkoutSvc}
}

func (h *checkoutHandler) Checkout(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	var items []entity.SKU
	err := json.NewDecoder(req.Body).Decode(&items)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	checkout := h.checkoutSvc.Checkout(req.Context(), items)
	b, _ := json.Marshal(checkout)
	rw.WriteHeader(http.StatusOK)
	rw.Header().Add("Content-Type", "Application/JSON")
	rw.Write(b)
}
