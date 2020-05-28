package item

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/data-portal-server/util/render"
	"github.com/go-chi/chi"
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

	orderID := chi.URLParam(r, "order_id")
	id, _ := strconv.Atoi(orderID)

	result := paging{}

	offset, limit := util.Paging(r.URL.Query())

	model.DB.Preload("Product").Preload("Product.ProductType").Preload("Product.Currency").Model(&model.OrderItem{}).Where(&model.OrderItem{OrderID: uint(id)}).Count(&result.Total).Offset(offset).Limit(limit).Find(&result.Nodes)

	render.JSON(w, http.StatusOK, result)
}
