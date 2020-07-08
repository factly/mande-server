package payment

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/renderx"
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
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &model.Payment{}

	result.ID = uint(id)

	// check record exists or not
	err = model.DB.First(&result).Error
	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	model.DB.Delete(&result)

	renderx.JSON(w, http.StatusOK, nil)
}
