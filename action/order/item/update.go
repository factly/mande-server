package item

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
)

// updateOrderItem - Update orderItems by id
// @Summary Update a orderItems by id
// @Description Update orderItems by ID
// @Tags OrderItem
// @ID update-orderItems-by-id
// @Produce json
// @Consume json
// @Param id path string true "Order ID"
// @Param oid path string true "OrderItem ID"
// @Param OrderItem body orderItem false "OrderItem"
// @Success 200 {object} model.OrderItem
// @Failure 400 {array} string
// @Router /orders/{id}/orderItems/{oid} [put]
func updateOrderItem(w http.ResponseWriter, r *http.Request) {

	orderItemID := chi.URLParam(r, "oid")
	id, err := strconv.Atoi(orderItemID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &model.OrderItem{}
	orderItem := &model.OrderItem{
		ID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&req)

	model.DB.Model(&orderItem).Updates(model.OrderItem{
		ExtraInfo: req.ExtraInfo,
		ProductID: req.ProductID,
		OrderID:   req.OrderID,
	})
	model.DB.Preload("Product").Preload("Product.Status").Preload("Product.ProductType").Preload("Product.Currency").Preload("Order").First(&orderItem)

	json.NewEncoder(w).Encode(orderItem)
}
