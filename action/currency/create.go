package currency

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
)

// create - Create currency
// @Summary Create currency
// @Description Create currency
// @Tags Currency
// @ID add-currency
// @Consume json
// @Produce  json
// @Param Currency body currency true "Currency object"
// @Success 201 {object} model.Currency
// @Failure 400 {array} string
// @Router /currencies [post]
func create(w http.ResponseWriter, r *http.Request) {

	currency := &currency{}

	json.NewDecoder(r.Body).Decode(&currency)

	err := validation.Validator.Struct(currency)
	if err != nil {
		validation.ValidatorErrors(w, r, err)
		return
	}

	result := &model.Currency{
		Name:    currency.Name,
		IsoCode: currency.IsoCode,
	}

	err = model.DB.Model(&model.Currency{}).Create(&result).Error

	if err != nil {
		log.Fatal(err)
	}

	render.JSON(w, http.StatusCreated, result)
}
