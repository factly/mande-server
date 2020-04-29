package order

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
)

// updateOrder - Update orders by id
// @Summary Update a orders by id
// @Description Update orders by ID
// @Tags Order
// @ID update-orders-by-id
// @Produce json
// @Consume json
// @Param id path string true "Order ID"
// @Param Order body order false "Order"
// @Success 200 {object} model.Order
// @Failure 400 {array} string
// @Router /orders/{id} [put]
func updateOrder(w http.ResponseWriter, r *http.Request) {

	orderID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(orderID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &model.Order{}
	orders := &model.Order{
		ID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&req)

	model.DB.Model(&orders).Updates(model.Order{
		UserID:    req.UserID,
		PaymentID: req.PaymentID,
		Status:    req.Status,
		CartID:    req.CartID,
	})
	model.DB.Preload("Payment").Preload("Payment.Currency").Preload("Cart").First(&orders)

	json.NewEncoder(w).Encode(orders)
}
