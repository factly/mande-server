package product

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// update - Update product by id
// @Summary Update a product by id
// @Description Update product by ID
// @Tags Product
// @ID update-product-by-id
// @Produce json
// @Consume json
// @Param product_id path string true "Product ID"
// @Param Product body product false "Product"
// @Success 200 {object} model.Product
// @Failure 400 {array} string
// @Router /products/{product_id} [put]
func update(w http.ResponseWriter, r *http.Request) {

	productID := chi.URLParam(r, "product_id")
	id, err := strconv.Atoi(productID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	product := &product{}
	categories := []model.ProductCategory{}
	tags := []model.ProductTag{}
	json.NewDecoder(r.Body).Decode(&product)

	result := &productData{}
	result.ID = uint(id)

	model.DB.Model(&result.Product).Updates(&model.Product{
		CurrencyID:    product.CurrencyID,
		ProductTypeID: product.ProductTypeID,
		Status:        product.Status,
		Title:         product.Title,
		Price:         product.Price,
		Slug:          product.Slug,
	}).Preload("ProductType").Preload("Status").Preload("Currency").First(&result.Product)

	// fetch all categories
	model.DB.Model(&model.ProductCategory{}).Where(&model.ProductCategory{
		ProductID: uint(id),
	}).Preload("Category").Find(&categories)

	// fetch all tags
	model.DB.Model(&model.ProductTag{}).Where(&model.ProductTag{
		ProductID: uint(id),
	}).Preload("Tag").Find(&tags)

	// delete tags
	for _, t := range tags {
		present := false
		for _, id := range product.TagIDs {
			if t.TagID == id {
				present = true
			}
		}
		if present == false {
			model.DB.Where(&model.ProductTag{
				TagID:     t.TagID,
				ProductID: uint(id),
			}).Delete(model.ProductTag{})
		}
	}

	// creating new tags
	for _, id := range product.TagIDs {
		present := false
		for _, t := range tags {
			if t.TagID == id {
				present = true
				result.Tags = append(result.Tags, t.Tag)
			}
		}
		if present == false {
			productTag := &model.ProductTag{}
			productTag.TagID = uint(id)
			productTag.ProductID = result.ID

			err = model.DB.Model(&model.ProductTag{}).Create(&productTag).Error

			if err != nil {
				return
			}
			model.DB.Model(&model.ProductTag{}).Preload("Tag").First(&productTag)
			result.Tags = append(result.Tags, productTag.Tag)
		}
	}

	// delete categories
	for _, c := range categories {
		present := false
		for _, id := range product.CategoryIDs {
			if c.CategoryID == id {
				present = true
			}
		}
		if present == false {
			model.DB.Where(&model.ProductCategory{
				CategoryID: c.CategoryID,
				ProductID:  uint(id),
			}).Delete(model.ProductCategory{})
		}
	}

	// creating new categories
	for _, id := range product.CategoryIDs {
		present := false
		for _, c := range categories {
			if c.CategoryID == id {
				present = true
				result.Categories = append(result.Categories, c.Category)
			}
		}
		if present == false {
			productCategory := &model.ProductCategory{}
			productCategory.CategoryID = uint(id)
			productCategory.ProductID = result.ID

			err = model.DB.Model(&model.ProductCategory{}).Create(&productCategory).Error

			if err != nil {
				return
			}

			model.DB.Model(&model.ProductCategory{}).Preload("Category").First(&productCategory)
			result.Categories = append(result.Categories, productCategory.Category)
		}
	}

	renderx.JSON(w, http.StatusOK, result)
}
