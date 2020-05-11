package payment

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// delete - Delete payment by id
// @Summary Delete a payment
// @Description Delete payment by ID
// @Tags Payment
// @ID delete-payment-by-id
// @Consume  json
// @Param payment_id path string true "Payment ID"
// @Success 200
// @Failure 400 {array} string
// @Router /payments/{payment_id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {

	paymentID := chi.URLParam(r, "payment_id")
	id, err := strconv.Atoi(paymentID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	payment := &model.Payment{}

	payment.ID = uint(id)

	// check record exists or not
	err = model.DB.First(&payment).Error
	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	model.DB.Delete(&payment)

	render.JSON(w, http.StatusOK, nil)
}
