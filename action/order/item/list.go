package item

import (
	"encoding/json"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
)

// list - Get all order items
// @Summary Show all order items
// @Description Get all order items
// @Tags OrderItem
// @ID get-all-order-items
// @Produce  json
// @Param order_id path string true "Order ID"
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} model.OrderItem
// @Router /orders/{order_id}/items [get]
func list(w http.ResponseWriter, r *http.Request) {

	var orderItems []model.OrderItem

	offset, limit := util.Paging(r.URL.Query())

	model.DB.Offset(offset).Limit(limit).Preload("Product").Preload("Product.Status").Preload("Product.ProductType").Preload("Product.Currency").Preload("Order").Model(&model.OrderItem{}).Find(&orderItems)

	json.NewEncoder(w).Encode(orderItems)
}
