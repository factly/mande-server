package order

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/data-portal-server/validation"
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
		validation.InvalidID(w, r)
		return
	}

	req := &model.Order{}
	req.ID = uint(id)

	err = model.DB.Model(&model.Order{}).Preload("Payment").Preload("Payment.Currency").Preload("Cart").First(&req).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	util.Render(w, http.StatusOK, req)
}
