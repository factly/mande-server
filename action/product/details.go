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
// @Param id path string true "Product ID"
// @Success 200 {object} model.Product
// @Failure 400 {array} string
// @Router /products/{id} [get]
func details(w http.ResponseWriter, r *http.Request) {

	productID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(productID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	result := &model.Product{}

	result.ID = uint(id)

	err = model.DB.Model(&model.Product{}).First(&result).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	model.DB.Preload("ProductType").Preload("Status").Preload("Currency").First(&result)

	render.JSON(w, http.StatusOK, result)
}
