package payment

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
)

// getPayments - Get all payments
// @Summary Show all payments
// @Description Get all payments
// @Tags Payment
// @ID get-all-payments
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} model.Payment
// @Router /payments [get]
func getPayments(w http.ResponseWriter, r *http.Request) {

	var payments []model.Payment
	p := r.URL.Query().Get("page")
	pg, _ := strconv.Atoi(p) // pg contains page number
	l := r.URL.Query().Get("limit")
	li, _ := strconv.Atoi(l) // li contains perPage number

	offset := 0 // no. of records to skip
	limit := 5  // limt

	if li > 0 && li <= 10 {
		limit = li
	}

	if pg > 1 {
		offset = (pg - 1) * limit
	}

	model.DB.Offset(offset).Limit(limit).Preload("Currency").Model(&model.Payment{}).Find(&payments)

	json.NewEncoder(w).Encode(payments)
}