package item

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// update - Update orderItems by id
// @Summary Update a orderItems by id
// @Description Update orderItems by ID
// @Tags OrderItem
// @ID update-orderItems-by-id
// @Produce json
// @Consume json
// @Param order_id path string true "Order ID"
// @Param item_id path string true "OrderItem ID"
// @Param OrderItem body orderItem false "OrderItem"
// @Success 200 {object} model.OrderItem
// @Failure 400 {array} string
// @Router /orders/{order_id}/items/{item_id} [put]
func update(w http.ResponseWriter, r *http.Request) {

	orderID := chi.URLParam(r, "order_id")
	oid, err := strconv.Atoi(orderID)

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	orderItemID := chi.URLParam(r, "item_id")
	id, err := strconv.Atoi(orderItemID)

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	orderItem := &orderItem{}

	json.NewDecoder(r.Body).Decode(&orderItem)

	result := &model.OrderItem{}
	result.ID = uint(id)
	result.OrderID = uint(oid)

	model.DB.Model(&result).Updates(model.OrderItem{
		ExtraInfo: orderItem.ExtraInfo,
		ProductID: orderItem.ProductID,
	}).Preload("Product").Preload("Product.Currency").First(&result)

	renderx.JSON(w, http.StatusOK, result)
}
