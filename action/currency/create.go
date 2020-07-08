package currency

import (
	"encoding/json"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
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

	validationError := validationx.Check(currency)
	if validationError != nil {
		errorx.Render(w, validationError)
		return
	}

	result := &model.Currency{
		Name:    currency.Name,
		IsoCode: currency.IsoCode,
	}

	err := model.DB.Model(&model.Currency{}).Create(&result).Error

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	renderx.JSON(w, http.StatusCreated, result)
}
