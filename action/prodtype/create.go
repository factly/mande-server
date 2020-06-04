package prodtype

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
)

// create - Create product type
// @Summary Create product type
// @Description Create product type
// @Tags Type
// @ID add-type
// @Consume json
// @Produce  json
// @Param Type body productType true "Type object"
// @Success 201 {object} model.ProductType
// @Failure 400 {array} string
// @Router /types [post]
func create(w http.ResponseWriter, r *http.Request) {

	productType := &productType{}
	json.NewDecoder(r.Body).Decode(&productType)

	err := validation.Validator.Struct(productType)
	if err != nil {
		validation.ValidatorErrors(w, r, err)
		return
	}

	result := &model.ProductType{
		Name: productType.Name,
	}

	err = model.DB.Model(&model.ProductType{}).Create(&result).Error

	if err != nil {
		log.Fatal(err)
	}

	render.JSON(w, http.StatusCreated, result)
}
