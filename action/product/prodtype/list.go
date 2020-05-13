package prodtype

import (
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/data-portal-server/util/render"
)

// list response
type paging struct {
	Total int                 `json:"total"`
	Nodes []model.ProductType `json:"nodes"`
}

// list - Get all productTypes
// @Summary Show all productTypes
// @Description Get all productTypes
// @Tags Type
// @ID get-all-productTypes
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /products/{product_id}/type [get]
func list(w http.ResponseWriter, r *http.Request) {

	data := paging{}

	offset, limit := util.Paging(r.URL.Query())

	model.DB.Model(&model.ProductType{}).Count(&data.Total).Offset(offset).Limit(limit).Find(&data.Nodes)

	render.JSON(w, http.StatusOK, data)
}
