package payment

import (
	"encoding/json"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
)

// list response
type paging struct {
	Total int             `json:"total"`
	Nodes []model.Payment `json:"nodes"`
}

// list - Get all payments
// @Summary Show all payments
// @Description Get all payments
// @Tags Payment
// @ID get-all-payments
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /payments [get]
func list(w http.ResponseWriter, r *http.Request) {

	data := paging{}

	offset, limit := util.Paging(r.URL.Query())

	model.DB.Offset(offset).Limit(limit).Preload("Currency").Model(&model.Payment{}).Find(&data.Nodes).Offset(0).Limit(-1).Count(&data.Total)

	json.NewEncoder(w).Encode(data)
}
