package item

import (
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/data-portal-server/util/render"
)

// list response
type paging struct {
	Total int               `json:"total"`
	Nodes []model.OrderItem `json:"nodes"`
}

// list - Get all order items
// @Summary Show all order items
// @Description Get all order items
// @Tags OrderItem
// @ID get-all-order-items
// @Produce  json
// @Param order_id path string true "Order ID"
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /orders/{order_id}/items [get]
func list(w http.ResponseWriter, r *http.Request) {

	data := paging{}

	offset, limit := util.Paging(r.URL.Query())

	model.DB.Preload("Product").Preload("Product.Status").Preload("Product.ProductType").Preload("Product.Currency").Preload("Order").Model(&model.OrderItem{}).Offset(offset).Limit(limit).Count(&data.Total).Find(&data.Nodes)

	render.JSON(w, http.StatusOK, data)
}
