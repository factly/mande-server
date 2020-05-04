package currency

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// delete - Delete currency by id
// @Summary Delete a currency
// @Description Delete currency by ID
// @Tags Currency
// @ID delete-currency-by-id
// @Consume  json
// @Param id path string true "Currency ID"
// @Success 200
// @Failure 400 {array} string
// @Router /currencies/{id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {
	currencyID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(currencyID)
	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	currency := &model.Currency{}
	currency.ID = uint(id)

	// check record exists or not
	err = model.DB.First(&currency).Error
	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	model.DB.Delete(&currency)

	render.JSON(w, http.StatusOK, nil)
}
