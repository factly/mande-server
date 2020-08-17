package currency

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// delete - Delete currency by id
// @Summary Delete a currency
// @Description Delete currency by ID
// @Tags Currency
// @ID delete-currency-by-id
// @Consume  json
// @Param currency_id path string true "Currency ID"
// @Success 200
// @Failure 400 {array} string
// @Router /currencies/{currency_id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {
	currencyID := chi.URLParam(r, "currency_id")
	id, err := strconv.Atoi(currencyID)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &model.Currency{}
	result.ID = uint(id)

	// check record exists or not
	err = model.DB.First(&result).Error
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	// check if currency is associated with payment
	var totAssociated int
	model.DB.Model(&model.Payment{}).Where(&model.Payment{
		CurrencyID: uint(id),
	}).Count(&totAssociated)

	if totAssociated != 0 {
		loggerx.Error(errors.New("currency is associated with payment"))
		errorx.Render(w, errorx.Parser(errorx.CannotSaveChanges()))
		return
	}

	// check if currency is associated with product
	model.DB.Model(&model.Product{}).Where(&model.Product{
		CurrencyID: uint(id),
	}).Count(&totAssociated)

	if totAssociated != 0 {
		loggerx.Error(errors.New("currency is associated with product"))
		errorx.Render(w, errorx.Parser(errorx.CannotSaveChanges()))
		return
	}

	// check if currency is associated with dataset
	model.DB.Model(&model.Dataset{}).Where(&model.Dataset{
		CurrencyID: uint(id),
	}).Count(&totAssociated)

	if totAssociated != 0 {
		loggerx.Error(errors.New("currency is associated with dataset"))
		errorx.Render(w, errorx.Parser(errorx.CannotSaveChanges()))
		return
	}

	model.DB.Delete(&result)

	renderx.JSON(w, http.StatusOK, nil)
}
