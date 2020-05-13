package currency

import (
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/data-portal-server/util/render"
)

// list response
type paging struct {
	Total int              `json:"total"`
	Nodes []model.Currency `json:"nodes"`
}

// list - Get all currencies
// @Summary Show all currencies
// @Description Get all currencies
// @Tags Currency
// @ID get-all-currencies
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /currencies [get]
func list(w http.ResponseWriter, r *http.Request) {

	result := paging{}

	offset, limit := util.Paging(r.URL.Query())

	model.DB.Model(&model.Currency{}).Count(&result.Total).Offset(offset).Limit(limit).Find(&result.Nodes)

	render.JSON(w, http.StatusOK, result)
}
