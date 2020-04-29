package item

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
)

// getOrderItems - Get all order items
// @Summary Show all order items
// @Description Get all order items
// @Tags OrderItem
// @ID get-all-order-items
// @Produce  json
// @Param id path string true "Order ID"
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} model.OrderItem
// @Router /orders/{id}/order-items [get]
func getOrderItems(w http.ResponseWriter, r *http.Request) {

	var orderItems []model.OrderItem
	p := r.URL.Query().Get("page")
	pg, _ := strconv.Atoi(p) // pg contains page number
	l := r.URL.Query().Get("limit")
	li, _ := strconv.Atoi(l) // li contains perPage number

	offset := 0 // no. of records to skip
	limit := 5  // limt

	if li > 0 && li <= 10 {
		limit = li
	}

	if pg > 1 {
		offset = (pg - 1) * limit
	}

	model.DB.Offset(offset).Limit(limit).Preload("Product").Preload("Product.Status").Preload("Product.ProductType").Preload("Product.Currency").Preload("Order").Model(&model.OrderItem{}).Find(&orderItems)

	json.NewEncoder(w).Encode(orderItems)
}
