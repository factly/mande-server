package cart

import (
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
)

// list response
type paging struct {
	Total int64        `json:"total"`
	Nodes []model.Cart `json:"nodes"`
}

// list - Get all carts
// @Summary Show all carts
// @Description Get all carts
// @Tags Cart
// @ID get-all-carts
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /carts [get]
func list(w http.ResponseWriter, r *http.Request) {

	result := paging{}
	result.Nodes = make([]model.Cart, 0)

	offset, limit := paginationx.Parse(r.URL.Query())

	model.DB.Model(&model.Cart{}).Preload("Products").Preload("Products.Tags").Preload("Products.Datasets").Preload("Products.Currency").Count(&result.Total).Offset(offset).Limit(limit).Find(&result.Nodes)

	renderx.JSON(w, http.StatusOK, result)
}
