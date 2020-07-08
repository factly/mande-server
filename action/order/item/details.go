package item

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/renderx"
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

	result := &model.OrderItem{}
	result.ID = uint(id)
	result.OrderID = uint(oid)

	err = model.DB.Model(&model.OrderItem{}).Preload("Product").Preload("Product.Currency").First(&result).Error

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	renderx.JSON(w, http.StatusOK, result)
}
