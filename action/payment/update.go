package payment

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// update - Update payment by id
// @Summary Update a payment by id
// @Description Update payment by ID
// @Tags Payment
// @ID update-payment-by-id
// @Produce json
// @Consume json
// @Param id path string true "Payment ID"
// @Param Payment body payment false "Payment"
// @Success 200 {object} model.Payment
// @Failure 400 {array} string
// @Router /payments/{id} [put]
func update(w http.ResponseWriter, r *http.Request) {

	paymentID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(paymentID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	payment := &model.Payment{
		ID: uint(id),
	}

	req := &model.Payment{}

	json.NewDecoder(r.Body).Decode(&req)

	model.DB.Model(&payment).Updates(&model.Payment{
		Amount:     req.Amount,
		Gateway:    req.Gateway,
		Status:     req.Status,
		CurrencyID: req.CurrencyID,
	})
	model.DB.First(&payment)
	model.DB.Model(&payment).Association("Currency").Find(&payment.Currency)
	json.NewEncoder(w).Encode(payment)
}
