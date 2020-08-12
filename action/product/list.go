package product

import (
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
)

// list response
type paging struct {
	Total int             `json:"total"`
	Nodes []model.Product `json:"nodes"`
}

// list - Get all products
// @Summary Show all products
// @Description Get all products
// @Tags Product
// @ID get-all-products
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /products [get]
func list(w http.ResponseWriter, r *http.Request) {

	result := &paging{}
	result.Nodes = make([]model.Product, 0)

	offset, limit := paginationx.Parse(r.URL.Query())

	model.DB.Model(&model.Product{}).Preload("Currency").Preload("FeaturedMedium").Preload("Tags").Preload("Datasets").Count(&result.Total).Offset(offset).Limit(limit).Find(&result.Nodes)

	renderx.JSON(w, http.StatusOK, result)
}
