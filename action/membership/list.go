package membership

import (
	"encoding/json"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
)

// list response
type paging struct {
	Total int                `json:"total"`
	Nodes []model.Membership `json:"nodes"`
}

// list - Get all memberships
// @Summary Show all memberships
// @Description Get all memberships
// @Tags Membership
// @ID get-all-memberships
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /memberships [get]
func list(w http.ResponseWriter, r *http.Request) {

	data := paging{}

	offset, limit := util.Paging(r.URL.Query())

	model.DB.Preload("User").Preload("Plan").Preload("Payment").Preload("Payment.Currency").Model(&model.Membership{}).Count(&data.Total).Offset(offset).Limit(limit).Find(&data.Nodes)

	json.NewEncoder(w).Encode(data)
}
