package currency

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
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
		validation.InvalidID(w, r)
		return
	}

	result := &model.Currency{}
	result.ID = uint(id)

	// check record exists or not
	err = model.DB.First(&result).Error
	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	model.DB.Delete(&result)

	renderx.JSON(w, http.StatusOK, nil)
}
