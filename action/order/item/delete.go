package item

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// delete - Delete order item by id
// @Summary Delete a order item
// @Description Delete order item by ID
// @Tags OrderItem
// @ID delete-order-items-by-id
// @Consume  json
// @Param item_id path string true "OrderItem ID"
// @Param order_id path string true "Order ID"
// @Success 200 {object} model.OrderItem
// @Failure 400 {array} string
// @Router /orders/{order_id}/items/{item_id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {

	orderItemID := chi.URLParam(r, "item_id")
	id, err := strconv.Atoi(orderItemID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	orderItems := &model.OrderItem{
		ID: uint(id),
	}

	// check record exists or not
	err = model.DB.First(&orderItems).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	model.DB.Preload("Product").Preload("Product.Status").Preload("Product.ProductType").Preload("Product.Currency").Preload("Order").Delete(&orderItems)

	json.NewEncoder(w).Encode(orderItems)
}
