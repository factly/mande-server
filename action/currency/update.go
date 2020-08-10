package currency

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// update - Update currency by id
// @Summary Update a currency by id
// @Description Update currency by ID
// @Tags Currency
// @ID update-currency-by-id
// @Produce json
// @Consume json
// @Param currency_id path string true "Currency ID"
// @Param Currency body currency false "Currency"
// @Success 200 {object} model.Currency
// @Failure 400 {array} string
// @Router /currencies/{currency_id} [put]
func update(w http.ResponseWriter, r *http.Request) {
	currencyID := chi.URLParam(r, "currency_id")
	id, err := strconv.Atoi(currencyID)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	currency := &currency{}

	json.NewDecoder(r.Body).Decode(&currency)

	result := &model.Currency{}
	result.ID = uint(id)

	model.DB.Model(&result).Updates(model.Currency{
		IsoCode: currency.IsoCode,
		Name:    currency.Name,
	}).First(&result)

	renderx.JSON(w, http.StatusOK, result)
}
