package product

import (
	"encoding/json"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
)

// list response
type paging struct {
	Total int             `json:"total"`
	Nodes []model.Product `json:"nodes"`
}

// list - Get all products
// @Summary Show all products
// @Description Get all products
// @Tags Product
// @ID get-all-products
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /products [get]
func list(w http.ResponseWriter, r *http.Request) {

	data := paging{}

	offset, limit := util.Paging(r.URL.Query())

	model.DB.Preload("Currency").Preload("Status").Preload("ProductType").Model(&model.Product{}).Count(&data.Total).Offset(offset).Limit(limit).Find(&data.Nodes)

	json.NewEncoder(w).Encode(data)
}
