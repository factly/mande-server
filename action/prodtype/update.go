package prodtype

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// update - Update product type by id
// @Summary Update a product type by id
// @Description Update product type by ID
// @Tags Type
// @ID update-product-type-by-id
// @Produce json
// @Consume json
// @Param type_id path string true "Product type ID"
// @Param type body productType false "Product type"
// @Success 200 {object} model.ProductType
// @Failure 400 {array} string
// @Router /types/{type_id} [put]
func update(w http.ResponseWriter, r *http.Request) {

	productTypeID := chi.URLParam(r, "type_id")
	id, err := strconv.Atoi(productTypeID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	productType := &productType{}
	result := &model.ProductType{}
	result.ID = uint(id)

	json.NewDecoder(r.Body).Decode(&productType)

	model.DB.Model(&result).Updates(model.ProductType{
		Name: productType.Name,
	}).First(&result)

	renderx.JSON(w, http.StatusOK, result)
}
