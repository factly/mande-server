package currency

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
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
// @Param id path string true "Currecny ID"
// @Param Currency body currency false "Currency"
// @Success 200 {object} model.Currency
// @Failure 400 {array} string
// @Router /currencies/{id} [put]
func update(w http.ResponseWriter, r *http.Request) {
	currencyID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(currencyID)
	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	currency := &model.Currency{}
	currency.ID = uint(id)
	req := &model.Currency{}

	json.NewDecoder(r.Body).Decode(&req)
	model.DB.Model(&currency).Updates(model.Currency{IsoCode: req.IsoCode, Name: req.Name})
	model.DB.First(&currency)

	util.Render(w, http.StatusOK, currency)
}
