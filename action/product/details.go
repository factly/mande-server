package product

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// details - Get product by id
// @Summary Show a product by id
// @Description Get product by ID
// @Tags Product
// @ID get-product-by-id
// @Produce  json
// @Param product_id path string true "Product ID"
// @Success 200 {object} model.Product
// @Failure 400 {array} string
// @Router /products/{product_id} [get]
func details(w http.ResponseWriter, r *http.Request) {

	var datasets []model.ProductDataset
	var tags []model.ProductTag

	productID := chi.URLParam(r, "product_id")
	id, err := strconv.Atoi(productID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	result := &productData{}
	result.Tags = make([]model.Tag, 0)
	result.Datasets = make([]model.Dataset, 0)

	result.Product.ID = uint(id)

	err = model.DB.Model(&model.Product{}).First(&result.Product).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	model.DB.Preload("Currency").Preload("FeaturedMedia").First(&result.Product)

	model.DB.Model(&model.ProductDataset{}).Where(&model.ProductDataset{
		ProductID: uint(id),
	}).Preload("Dataset").Find(&datasets)

	model.DB.Model(&model.ProductTag{}).Where(&model.ProductTag{
		ProductID: uint(id),
	}).Preload("Tag").Find(&tags)

	for _, d := range datasets {
		result.Datasets = append(result.Datasets, d.Dataset)
	}

	for _, t := range tags {
		result.Tags = append(result.Tags, t.Tag)
	}

	renderx.JSON(w, http.StatusOK, result)
}
