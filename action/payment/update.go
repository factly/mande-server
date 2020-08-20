package payment

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
	"github.com/go-chi/chi"
)

// update - Update payment by id
// @Summary Update a payment by id
// @Description Update payment by ID
// @Tags Payment
// @ID update-payment-by-id
// @Produce json
// @Consume json
// @Param payment_id path string true "Payment ID"
// @Param Payment body payment false "Payment"
// @Success 200 {object} model.Payment
// @Failure 400 {array} string
// @Router /payments/{payment_id} [put]
func update(w http.ResponseWriter, r *http.Request) {

	paymentID := chi.URLParam(r, "payment_id")
	id, err := strconv.Atoi(paymentID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	payment := &payment{}

	err = json.NewDecoder(r.Body).Decode(&payment)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}
	validationError := validationx.Check(payment)
	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	result := &model.Payment{}
	result.ID = uint(id)

	model.DB.Model(&result).Updates(&model.Payment{
		Amount:     payment.Amount,
		Gateway:    payment.Gateway,
		Status:     payment.Status,
		CurrencyID: payment.CurrencyID,
	}).First(&result).Preload("Currency").Find(&result.Currency)

	renderx.JSON(w, http.StatusOK, result)
}
