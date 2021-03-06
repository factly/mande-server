package plan

import (
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
)

// list response
type paging struct {
	Total int64        `json:"total"`
	Nodes []model.Plan `json:"nodes"`
}

// list - Get all plans
// @Summary Show all plans
// @Description Get all plans
// @Tags Plan
// @ID get-all-plans
// @Produce  json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /plans [get]
func list(w http.ResponseWriter, r *http.Request) {

	result := paging{}
	result.Nodes = make([]model.Plan, 0)

	offset, limit := paginationx.Parse(r.URL.Query())

	model.DB.Model(&model.Plan{}).Preload("Currency").Preload("Catalogs").Preload("Catalogs.Products").Preload("Catalogs.Products.Currency").Preload("Catalogs.Products.Datasets").Preload("Catalogs.Products.Tags").Count(&result.Total).Offset(offset).Limit(limit).Find(&result.Nodes)

	renderx.JSON(w, http.StatusOK, result)
}
