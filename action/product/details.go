package product

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
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

	var categories []model.ProductCategory
	var tags []model.ProductTag

	productID := chi.URLParam(r, "product_id")
	id, err := strconv.Atoi(productID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	result := &productData{}

	result.Product.ID = uint(id)

	err = model.DB.Model(&model.Product{}).First(&result.Product).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	model.DB.Preload("ProductType").Preload("Currency").First(&result.Product)

	model.DB.Model(&model.ProductCategory{}).Where(&model.ProductCategory{
		ProductID: uint(id),
	}).Preload("Category").Find(&categories)

	model.DB.Model(&model.ProductTag{}).Where(&model.ProductTag{
		ProductID: uint(id),
	}).Preload("Tag").Find(&tags)

	for _, c := range categories {
		result.Categories = append(result.Categories, c.Category)
	}

	for _, t := range tags {
		result.Tags = append(result.Tags, t.Tag)
	}

	render.JSON(w, http.StatusOK, result)
}
