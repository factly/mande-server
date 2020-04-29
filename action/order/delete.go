package order

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
)

// deleteOrder - Delete orders by id
// @Summary Delete a orders
// @Description Delete orders by ID
// @Tags Order
// @ID delete-orders-by-id
// @Consume  json
// @Param id path string true "Order ID"
// @Success 200 {object} model.Order
// @Failure 400 {array} string
// @Router /orders/{id} [delete]
func deleteOrder(w http.ResponseWriter, r *http.Request) {

	orderID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(orderID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	orders := &model.Order{
		ID: uint(id),
	}

	// check record exists or not
	err = model.DB.First(&orders).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	model.DB.Preload("Payment").Preload("Payment.Currency").Preload("Cart").Delete(&orders)

	json.NewEncoder(w).Encode(orders)
}
