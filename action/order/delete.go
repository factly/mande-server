package order

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// delete - Delete orders by id
// @Summary Delete a orders
// @Description Delete orders by ID
// @Tags Order
// @ID delete-orders-by-id
// @Consume  json
// @Param order_id path string true "Order ID"
// @Success 200
// @Failure 400 {array} string
// @Router /orders/{order_id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {

	orderID := chi.URLParam(r, "order_id")
	id, err := strconv.Atoi(orderID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &model.Order{}
	result.ID = uint(id)

	// check record exists or not
	err = model.DB.First(&result).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	model.DB.Where(&model.OrderItem{
		OrderID: uint(id),
	}).Delete(&model.OrderItem{})

	model.DB.Delete(&result)

	renderx.JSON(w, http.StatusOK, nil)
}
