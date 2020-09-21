package order

import (
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
)

// list response
type paging struct {
	Total int           `json:"total"`
	Nodes []model.Order `json:"nodes"`
}

// list - Get all orders
// @Summary Show all orders
// @Description Get all orders
// @Tags Order
// @ID get-all-orders
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /orders [get]
func list(w http.ResponseWriter, r *http.Request) {

	result := paging{}
	result.Nodes = make([]model.Order, 0)

	offset, limit := paginationx.Parse(r.URL.Query())

	model.DB.Preload("Payment").Preload("Payment.Currency").Preload("Products").Preload("Products.Datasets").Preload("Products.Tags").Model(&model.Order{}).Count(&result.Total).Offset(offset).Limit(limit).Find(&result.Nodes)

	renderx.JSON(w, http.StatusOK, result)
}
