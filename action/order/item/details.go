package item

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// details - Get order item by id
// @Summary Show a order item by id
// @Description Get order item by ID
// @Tags OrderItem
// @ID get-order-item-by-id
// @Produce  json
// @Param order_id path string true "Order ID"
// @Param item_id path string true "Order item ID"
// @Success 200 {object} model.OrderItem
// @Failure 400 {array} string
// @Router /orders/{order_id}/items/{item_id} [get]
func details(w http.ResponseWriter, r *http.Request) {

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

	err = model.DB.Model(&model.OrderItem{}).Preload("Product").Preload("Product.Status").Preload("Product.ProductType").Preload("Product.Currency").Preload("Order").First(&orderItem).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	render.JSON(w, http.StatusOK, orderItem)
}
