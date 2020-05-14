package payment

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// details - Get payment by id
// @Summary Show a payment by id
// @Description Get payment by ID
// @Tags Payment
// @ID get-payment-by-id
// @Produce  json
// @Param id path string true "Payment ID"
// @Success 200 {object} model.Payment
// @Failure 400 {array} string
// @Router /payments/{id} [get]
func details(w http.ResponseWriter, r *http.Request) {

	paymentID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(paymentID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	result := &model.Payment{}
	result.ID = uint(id)

	err = model.DB.Model(&model.Payment{}).First(&result).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	model.DB.Model(&result).Preload("Currency").Find(&result.Currency)

	render.JSON(w, http.StatusOK, result)
}
