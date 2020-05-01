package product

import (
	"encoding/json"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
)

// list - Get all products
// @Summary Show all products
// @Description Get all products
// @Tags Product
// @ID get-all-products
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} model.Product
// @Router /products [get]
func list(w http.ResponseWriter, r *http.Request) {

	var products []model.Product

	offset, limit := util.Paging(r.URL.Query())

	model.DB.Offset(offset).Limit(limit).Preload("Currency").Preload("Status").Preload("ProductType").Model(&model.Product{}).Find(&products)

	json.NewEncoder(w).Encode(products)
}
