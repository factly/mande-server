package order

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// details - Get orders by id
// @Summary Show a orders by id
// @Description Get orders by ID
// @Tags Order
// @ID get-orders-by-id
// @Produce  json
// @Param order_id path string true "Order ID"
// @Success 200 {object} model.Order
// @Failure 400 {array} string
// @Router /orders/{order_id} [get]
func details(w http.ResponseWriter, r *http.Request) {

	orderID := chi.URLParam(r, "order_id")
	id, err := strconv.Atoi(orderID)

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &model.Order{}
	result.ID = uint(id)

	err = model.DB.Model(&model.Order{}).Preload("Payment").Preload("Payment.Currency").Preload("Cart").First(&result).Error

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	renderx.JSON(w, http.StatusOK, result)
}
