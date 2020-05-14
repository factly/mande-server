package product

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-playground/validator/v10"
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

	validate := validator.New()
	err := validate.StructExcept(product, "ProductType", "Status", "Currency")
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	result := &model.Product{
		Title:         product.Title,
		Slug:          product.Slug,
		Price:         product.Price,
		ProductTypeID: product.ProductTypeID,
		StatusID:      product.StatusID,
		CurrencyID:    product.CurrencyID,
	}

	err = model.DB.Model(&model.Product{}).Create(&result).Error

	if err != nil {
		log.Fatal(err)
	}

	model.DB.Preload("ProductType").Preload("Status").Preload("Currency").First(&result)

	render.JSON(w, http.StatusCreated, result)
}
