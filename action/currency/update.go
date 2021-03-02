package currency

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/middlewarex"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
	"github.com/go-chi/chi"
)

// update - Update currency by id
// @Summary Update a currency by id
// @Description Update currency by ID
// @Tags Currency
// @ID update-currency-by-id
// @Produce json
// @Consume json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param currency_id path string true "Currency ID"
// @Param Currency body currency false "Currency"
// @Success 200 {object} model.Currency
// @Failure 400 {array} string
// @Router /currencies/{currency_id} [put]
func update(w http.ResponseWriter, r *http.Request) {
	uID, err := middlewarex.GetUser(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	currencyID := chi.URLParam(r, "currency_id")
	id, err := strconv.Atoi(currencyID)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	currency := &currency{}

	err = json.NewDecoder(r.Body).Decode(&currency)
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

	result := &model.Currency{}
	result.ID = uint(id)

	err = model.DB.First(&result).Error
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	model.DB.Model(&result).Updates(model.Currency{
		Base:    model.Base{UpdatedByID: uint(uID)},
		IsoCode: currency.IsoCode,
		Name:    currency.Name,
	}).First(&result)

	renderx.JSON(w, http.StatusOK, result)
}
