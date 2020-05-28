package prodtype

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-playground/validator/v10"
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

	validate := validator.New()
	err := validate.Struct(productType)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
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