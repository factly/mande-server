package currency

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// details - Get currency by id
// @Summary Show a currency by id
// @Description get currency by ID
// @Tags Currency
// @ID get-currency-by-id
// @Produce  json
// @Param currency_id path string false "Currency ID"
// @Success 200 {object} model.Currency
// @Failure 400 {array} string
// @Router /currencies/{currency_id} [get]
func details(w http.ResponseWriter, r *http.Request) {
	currencyID := chi.URLParam(r, "currency_id")
	id, err := strconv.Atoi(currencyID)
	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	currency := &model.Currency{}
	currency.ID = uint(id)

	err = model.DB.Model(&model.Currency{}).First(&currency).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	render.JSON(w, http.StatusOK, currency)
}
