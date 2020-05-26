package order

import (
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/data-portal-server/util/render"
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

	offset, limit := util.Paging(r.URL.Query())

	model.DB.Preload("Payment").Preload("Payment.Currency").Preload("Cart").Model(&model.Order{}).Count(&result.Total).Offset(offset).Limit(limit).Find(&result.Nodes)

	render.JSON(w, http.StatusOK, result)
}
