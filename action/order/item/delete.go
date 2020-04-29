package item

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
)

// deleteOrderItem - Delete order item by id
// @Summary Delete a order item
// @Description Delete order item by ID
// @Tags OrderItem
// @ID delete-order-items-by-id
// @Consume  json
// @Param oid path string true "OrderItem ID"
// @Param id path string true "Order ID"
// @Success 200 {object} model.OrderItem
// @Failure 400 {array} string
// @Router /orders/{id}/order-items/{oid} [delete]
func deleteOrderItem(w http.ResponseWriter, r *http.Request) {

	orderItemID := chi.URLParam(r, "oid")
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