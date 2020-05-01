package prodtype

import (
	"encoding/json"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
)

// list - Get all productTypes
// @Summary Show all productTypes
// @Description Get all productTypes
// @Tags Type
// @ID get-all-productTypes
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} model.ProductType
// @Router /products/{id}/type [get]
func list(w http.ResponseWriter, r *http.Request) {

	var productTypes []model.ProductType

	offset, limit := util.Paging(r.URL.Query())

	model.DB.Offset(offset).Limit(limit).Model(&model.ProductType{}).Find(&productTypes)

	json.NewEncoder(w).Encode(productTypes)
}
