package product

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// update - Update product by id
// @Summary Update a product by id
// @Description Update product by ID
// @Tags Product
// @ID update-product-by-id
// @Produce json
// @Consume json
// @Param id path string true "Product ID"
// @Param Product body product false "Product"
// @Success 200 {object} model.Product
// @Failure 400 {array} string
// @Router /products/{id} [put]
func update(w http.ResponseWriter, r *http.Request) {

	productID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(productID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	product := &product{}
	json.NewDecoder(r.Body).Decode(&product)

	result := &model.Product{}
	result.ID = uint(id)

	model.DB.Model(&result).Updates(&model.Product{
		CurrencyID:    product.CurrencyID,
		ProductTypeID: product.ProductTypeID,
		StatusID:      product.StatusID,
		Title:         product.Title,
		Price:         product.Price,
		Slug:          product.Slug,
	}).First(&result)

	model.DB.Model(&result).Association("ProductType").Find(&result.ProductType)
	model.DB.Model(&result).Association("Currency").Find(&result.Currency)
	model.DB.Model(&result).Association("Status").Find(&result.Status)

	render.JSON(w, http.StatusOK, result)
}
