package payment

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
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

	req := &model.Payment{
		ID: uint(id),
	}

	err = model.DB.Model(&model.Payment{}).First(&req).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	model.DB.Model(&req).Association("Currency").Find(&req.Currency)
	json.NewEncoder(w).Encode(req)
}
