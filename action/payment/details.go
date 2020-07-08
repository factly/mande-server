package payment

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// details - Get payment by id
// @Summary Show a payment by id
// @Description Get payment by ID
// @Tags Payment
// @ID get-payment-by-id
// @Produce  json
// @Param payment_id path string true "Payment ID"
// @Success 200 {object} model.Payment
// @Failure 400 {array} string
// @Router /payments/{payment_id} [get]
func details(w http.ResponseWriter, r *http.Request) {

	paymentID := chi.URLParam(r, "payment_id")
	id, err := strconv.Atoi(paymentID)

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &model.Payment{}
	result.ID = uint(id)

	err = model.DB.Model(&model.Payment{}).First(&result).Error

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	model.DB.Model(&result).Preload("Currency").Find(&result.Currency)

	renderx.JSON(w, http.StatusOK, result)
}
