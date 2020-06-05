package prodtype

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// delete - Delete product type by id
// @Summary Delete a product type
// @Description Delete product type by ID
// @Tags Type
// @ID delete-product-type-by-id
// @Consume  json
// @Param type_id path string true "Product Type ID"
// @Success 200
// @Failure 400 {array} string
// @Router /types/{type_id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {

	productTypeID := chi.URLParam(r, "type_id")
	id, err := strconv.Atoi(productTypeID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	result := &model.ProductType{}

	result.ID = uint(id)

	// check record exists or not
	err = model.DB.First(&result).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	model.DB.Delete(&result)

	renderx.JSON(w, http.StatusOK, nil)
}
