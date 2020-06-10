package format

import (
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
)

// list response
type paging struct {
	Total int            `json:"total"`
	Nodes []model.Format `json:"nodes"`
}

// list - Get all formats
// @Summary Show all formats
// @Description Get all formats
// @Tags Format
// @ID get-all-formats
// @Produce  json
// @Param limit query string false "limit per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /formats [get]
func list(w http.ResponseWriter, r *http.Request) {

	result := paging{}

	offset, limit := paginationx.Parse(r.URL.Query())

	model.DB.Model(&model.Format{}).Count(&result.Total).Offset(offset).Limit(limit).Find(&result.Nodes)

	renderx.JSON(w, http.StatusOK, result)
}
