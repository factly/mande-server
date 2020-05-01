package payment

import (
	"encoding/json"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
)

// list - Get all payments
// @Summary Show all payments
// @Description Get all payments
// @Tags Payment
// @ID get-all-payments
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} model.Payment
// @Router /payments [get]
func list(w http.ResponseWriter, r *http.Request) {

	var payments []model.Payment

	offset, limit := util.Paging(r.URL.Query())

	model.DB.Offset(offset).Limit(limit).Preload("Currency").Model(&model.Payment{}).Find(&payments)

	json.NewEncoder(w).Encode(payments)
}
