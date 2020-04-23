package actions

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-api/models"
	"github.com/factly/data-portal-api/validationerrors"
	"github.com/go-playground/validator/v10"
)

// productType request body
type productType struct {
	Name string `json:"name"`
}

// CreateProductType - Create product type
// @Summary Create product type
// @Description Create product type
// @Tags Type
// @ID add-type
// @Consume json
// @Produce  json
// @Param Type body productType true "Type object"
// @Success 200 {object} models.ProductType
// @Failure 400 {array} string
// @Router /products/{id}/type [post]
func CreateProductType(w http.ResponseWriter, r *http.Request) {

	req := &models.ProductType{}
	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		msg := err.Error()
		validationerrors.ValidErrors(w, r, msg)
		return
	}

	err = models.DB.Model(&models.ProductType{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(req)
}
