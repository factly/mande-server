package currency

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-playground/validator/v10"
)

// create - Create currency
// @Summary Create currency
// @Description Create currency
// @Tags Currency
// @ID add-currency
// @Consume json
// @Produce  json
// @Param Currency body currency true "Currency object"
// @Success 200 {object} model.Currency
// @Failure 400 {array} string
// @Router /currencies [post]
func create(w http.ResponseWriter, r *http.Request) {

	req := &model.Currency{}

	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = model.DB.Model(&model.Currency{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(req)
}