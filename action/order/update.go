package order

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// update - Update orders by id
// @Summary Update a orders by id
// @Description Update orders by ID
// @Tags Order
// @ID update-orders-by-id
// @Produce json
// @Consume json
// @Param order_id path string true "Order ID"
// @Param Order body order false "Order"
// @Success 200 {object} model.Order
// @Failure 400 {array} string
// @Router /orders/{order_id} [put]
func update(w http.ResponseWriter, r *http.Request) {

	orderID := chi.URLParam(r, "order_id")
	id, err := strconv.Atoi(orderID)

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	order := &order{}

	json.NewDecoder(r.Body).Decode(&order)

	result := &model.Order{}
	result.ID = uint(id)

	model.DB.Model(&result).Updates(model.Order{
		UserID:    order.UserID,
		PaymentID: order.PaymentID,
		Status:    order.Status,
		CartID:    order.CartID,
	}).Preload("Payment").Preload("Payment.Currency").Preload("Cart").First(&result)

	renderx.JSON(w, http.StatusOK, result)
}
