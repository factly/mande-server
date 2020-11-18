package currency

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
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

	err := json.NewDecoder(r.Body).Decode(&currency)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	validationError := validationx.Check(currency)
	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	// Check if currency already exists
	var totCurrencies int64
	model.DB.Model(&model.Currency{}).Count(&totCurrencies)

	if totCurrencies > 0 {
		errorx.Render(w, errorx.Parser(errorx.Message{
			Message: "Cannot add more than one currency",
			Code:    422,
		}))
		return
	}

	result := &model.Currency{
		Name:    currency.Name,
		IsoCode: currency.IsoCode,
	}

	err = model.DB.Model(&model.Currency{}).Create(&result).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	renderx.JSON(w, http.StatusCreated, result)
}
