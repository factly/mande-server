package item

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
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
	oid, err := strconv.Atoi(orderID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	orderItemID := chi.URLParam(r, "item_id")
	id, err := strconv.Atoi(orderItemID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &model.OrderItem{}
	result.ID = uint(id)
	result.OrderID = uint(oid)

	// check record exists or not
	err = model.DB.First(&result).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}
	model.DB.Delete(&result)

	renderx.JSON(w, http.StatusOK, nil)
}
