package cart

import (
	"encoding/json"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
)

// list response
type paging struct {
	Total int64        `json:"total"`
	Nodes []model.Cart `json:"nodes"`
}

// list - Get all carts
// @Summary Show all carts
// @Description Get all carts
// @Tags Cart
// @ID get-all-carts
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /carts [get]
func list(w http.ResponseWriter, r *http.Request) {

	data := paging{}

	offset, limit := util.Paging(r.URL.Query())

	model.DB.Offset(offset).Limit(limit).Model(&model.Cart{}).Find(&data.Nodes).Offset(0).Limit(-1).Count(&data.Total)

	json.NewEncoder(w).Encode(data)
}
