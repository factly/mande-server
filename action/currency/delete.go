package currency

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
)

// deleteCurrency - Delete currency by id
// @Summary Delete a currency
// @Description Delete currency by ID
// @Tags Currency
// @ID delete-currency-by-id
// @Consume  json
// @Param id path string true "Currency ID"
// @Success 200 {object} model.Currency
// @Failure 400 {array} string
// @Router /currencies/{id} [delete]
func deleteCurrency(w http.ResponseWriter, r *http.Request) {
	currencyID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(currencyID)
	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	currency := &model.Currency{
		ID: uint(id),
	}

	// check record exists or not
	err = model.DB.First(&currency).Error
	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	model.DB.Delete(&currency)

	json.NewEncoder(w).Encode(currency)
}
