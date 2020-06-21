package product

import (
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
)

// list response
type paging struct {
	Total int           `json:"total"`
	Nodes []productData `json:"nodes"`
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

	var nodes []productData
	var products []model.Product
	result := &paging{}

	offset, limit := paginationx.Parse(r.URL.Query())

	model.DB.Preload("Currency").Preload("FeaturedMedia").Model(&model.Product{}).Count(&result.Total).Offset(offset).Limit(limit).Find(&products)

	for _, product := range products {
		var datasets []model.ProductDataset
		var tags []model.ProductTag
		data := &productData{}

		model.DB.Model(&model.ProductDataset{}).Where(&model.ProductDataset{
			ProductID: uint(product.ID),
		}).Preload("Dataset").Find(&datasets)

		model.DB.Model(&model.ProductTag{}).Where(&model.ProductTag{
			ProductID: uint(product.ID),
		}).Preload("Tag").Find(&tags)

		for _, d := range datasets {
			data.Datasets = append(data.Datasets, d.Dataset)
		}

		for _, t := range tags {
			data.Tags = append(data.Tags, t.Tag)
		}

		data.Product = product

		nodes = append(nodes, *data)
	}
	result.Nodes = nodes

	renderx.JSON(w, http.StatusOK, result)
}
