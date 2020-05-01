package cart

import (
	"encoding/json"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
)

// list - Get all carts
// @Summary Show all carts
// @Description Get all carts
// @Tags Cart
// @ID get-all-carts
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} model.Cart
// @Router /carts [get]
func list(w http.ResponseWriter, r *http.Request) {

	var carts []model.Cart

	offset, limit := util.Paging(r.URL.Query())

	model.DB.Offset(offset).Limit(limit).Model(&model.Cart{}).Find(&carts)

	json.NewEncoder(w).Encode(carts)
}
