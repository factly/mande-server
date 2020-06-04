package product

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
)

// create - Create product
// @Summary Create product
// @Description Create product
// @Tags Product
// @ID add-product
// @Consume json
// @Produce  json
// @Param Product body product true "Product object"
// @Success 201 {object} model.Product
// @Failure 400 {array} string
// @Router /products [post]
func create(w http.ResponseWriter, r *http.Request) {

	product := &product{}
	json.NewDecoder(r.Body).Decode(&product)

	err := validation.Validator.Struct(product)
	if err != nil {
		validation.ValidatorErrors(w, r, err)
		return
	}
	result := &productData{}
	result.Product = model.Product{
		Title:         product.Title,
		Slug:          product.Slug,
		Price:         product.Price,
		ProductTypeID: product.ProductTypeID,
		Status:        product.Status,
		CurrencyID:    product.CurrencyID,
	}

	err = model.DB.Model(&model.Product{}).Create(&result.Product).Error

	if err != nil {
		log.Fatal(err)
	}

	model.DB.Preload("ProductType").Preload("Currency").First(&result.Product)

	for _, id := range product.CategoryIDs {
		productCategory := &model.ProductCategory{}

		productCategory.CategoryID = uint(id)
		productCategory.ProductID = result.ID

		err = model.DB.Model(&model.ProductCategory{}).Create(&productCategory).Error

		if err != nil {
			log.Fatal(err)
		}
		model.DB.Model(&model.ProductCategory{}).Preload("Category").First(&productCategory)
		result.Categories = append(result.Categories, productCategory.Category)
	}

	for _, id := range product.TagIDs {
		productTag := &model.ProductTag{}

		productTag.TagID = uint(id)
		productTag.ProductID = result.ID

		err = model.DB.Model(&model.ProductTag{}).Create(&productTag).Error

		if err != nil {
			log.Fatal(err)
		}
		model.DB.Model(&model.ProductTag{}).Preload("Tag").First(&productTag)
		result.Tags = append(result.Tags, productTag.Tag)
	}

	render.JSON(w, http.StatusCreated, result)
}
