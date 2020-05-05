package item

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
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
// @Success 200
// @Failure 400 {array} string
// @Router /orders/{order_id}/items/{item_id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {

	orderID := chi.URLParam(r, "order_id")
	oid, _ := strconv.Atoi(orderID)

	orderItemID := chi.URLParam(r, "item_id")
	id, err := strconv.Atoi(orderItemID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	orderItem := &model.OrderItem{}
	orderItem.ID = uint(id)
	orderItem.OrderID = uint(oid)

	// check record exists or not
	err = model.DB.First(&orderItem).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	model.DB.Delete(&orderItem)

	render.JSON(w, http.StatusOK, nil)
}
