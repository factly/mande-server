package membership

import (
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
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

	result := paging{}
	result.Nodes = make([]model.Membership, 0)

	offset, limit := paginationx.Parse(r.URL.Query())

	model.DB.Preload("Plan").Preload("Plan.Catalogs").Preload("Payment").Preload("Payment.Currency").Model(&model.Membership{}).Count(&result.Total).Offset(offset).Limit(limit).Find(&result.Nodes)

	renderx.JSON(w, http.StatusOK, result)
}
