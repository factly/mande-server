package prodtype

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
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
// @Success 200 {object} model.ProductType
// @Failure 400 {array} string
// @Router /products/{id}/type [post]
func create(w http.ResponseWriter, r *http.Request) {

	req := &model.ProductType{}
	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = model.DB.Model(&model.ProductType{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(req)
}
