package catalog

import (
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
)

// list response
type paging struct {
	Total int             `json:"total"`
	Nodes []model.Catalog `json:"nodes"`
}

// list - Get all catalogs
// @Summary Show all catalogs
// @Description Get all catalogs
// @Tags Catalog
// @ID get-all-catalogs
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /catalogs [get]
func list(w http.ResponseWriter, r *http.Request) {

	result := paging{}
	result.Nodes = make([]model.Catalog, 0)
	offset, limit := paginationx.Parse(r.URL.Query())

	model.DB.Preload("FeaturedMedium").Preload("Products").Preload("Products.Currency").Preload("Products.FeaturedMedium").Preload("Products.Tags").Preload("Products.Datasets").Model(&model.Catalog{}).Count(&result.Total).Offset(offset).Limit(limit).Find(&result.Nodes)

	renderx.JSON(w, http.StatusOK, result)
}
