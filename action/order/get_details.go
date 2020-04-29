package order

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
)

// getOrderByID - Get orders by id
// @Summary Show a orders by id
// @Description Get orders by ID
// @Tags Order
// @ID get-orders-by-id
// @Produce  json
// @Param id path string true "Order ID"
// @Success 200 {object} model.Order
// @Failure 400 {array} string
// @Router /orders/{id} [get]
func getOrderByID(w http.ResponseWriter, r *http.Request) {

	orderID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(orderID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &model.Order{
		ID: uint(id),
	}

	err = model.DB.Model(&model.Order{}).Preload("Payment").Preload("Payment.Currency").Preload("Cart").First(&req).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(req)
}
