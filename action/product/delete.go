package product

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// delete - Delete product by id
// @Summary Delete a product
// @Description Delete product by ID
// @Tags Product
// @ID delete-product-by-id
// @Consume  json
// @Param product_id path string true "Product ID"
// @Success 200
// @Failure 400 {array} string
// @Router /products/{product_id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {

	productID := chi.URLParam(r, "product_id")
	id, err := strconv.Atoi(productID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	result := &model.Product{}

	result.ID = uint(id)

	// check record exists or not
	err = model.DB.First(&result).Error
	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	model.DB.Delete(&result)

	render.JSON(w, http.StatusOK, nil)
}
