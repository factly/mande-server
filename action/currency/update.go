package currency

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
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
		validation.InvalidID(w, r)
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

	render.JSON(w, http.StatusOK, result)
}
