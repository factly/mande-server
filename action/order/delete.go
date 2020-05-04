package order

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// delete - Delete orders by id
// @Summary Delete a orders
// @Description Delete orders by ID
// @Tags Order
// @ID delete-orders-by-id
// @Consume  json
// @Param order_id path string true "Order ID"
// @Success 200 {object} model.Order
// @Failure 400 {array} string
// @Router /orders/{order_id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {

	orderID := chi.URLParam(r, "order_id")
	id, err := strconv.Atoi(orderID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	orders := &model.Order{}
	orders.ID = uint(id)

	// check record exists or not
	err = model.DB.First(&orders).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	model.DB.Preload("Payment").Preload("Payment.Currency").Preload("Cart").Delete(&orders)

	util.Render(w, http.StatusOK, orders)
}
