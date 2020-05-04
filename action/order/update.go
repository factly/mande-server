package order

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/data-portal-server/validation"
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
		validation.InvalidID(w, r)
		return
	}

	req := &model.Order{}
	orders := &model.Order{}
	orders.ID = uint(id)

	json.NewDecoder(r.Body).Decode(&req)

	model.DB.Model(&orders).Updates(model.Order{
		UserID:    req.UserID,
		PaymentID: req.PaymentID,
		Status:    req.Status,
		CartID:    req.CartID,
	})
	model.DB.Preload("Payment").Preload("Payment.Currency").Preload("Cart").First(&orders)

	util.Render(w, http.StatusOK, orders)
}
