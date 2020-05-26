package prodtype

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// details - Get product type by id
// @Summary Show a product type by id
// @Description Get product type by ID
// @Tags Type
// @ID get-product-type-by-id
// @Produce  json
// @Param type_id path string true "Type ID"
// @Success 200 {object} model.ProductType
// @Failure 400 {array} string
// @Router /types/{type_id} [get]
func details(w http.ResponseWriter, r *http.Request) {

	productTypeID := chi.URLParam(r, "type_id")
	id, err := strconv.Atoi(productTypeID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	result := &model.ProductType{}
	result.ID = uint(id)

	err = model.DB.Model(&model.ProductType{}).First(&result).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	render.JSON(w, http.StatusOK, result)
}
