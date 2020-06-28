package catalog

import (
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
)

// list response
type paging struct {
	Total int           `json:"total"`
	Nodes []catalogData `json:"nodes"`
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
	result.Nodes = make([]catalogData, 0)
	nodes := make([]catalogData, 0)
	catalogs := []model.Catalog{}

	offset, limit := paginationx.Parse(r.URL.Query())

	model.DB.Preload("FeaturedMedia").Model(&model.Catalog{}).Count(&result.Total).Offset(offset).Limit(limit).Find(&catalogs)

	for _, catalog := range catalogs {
		var products []model.CatalogProduct
		data := &catalogData{}
		data.Products = make([]model.Product, 0)

		model.DB.Model(&model.CatalogProduct{}).Where(&model.CatalogProduct{
			CatalogID: uint(catalog.ID),
		}).Preload("Product").Find(&products)

		for _, t := range products {
			data.Products = append(data.Products, t.Product)
		}

		data.Catalog = catalog

		nodes = append(nodes, *data)
	}
	result.Nodes = nodes

	renderx.JSON(w, http.StatusOK, result)
}
