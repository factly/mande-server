package currency

import (
	"encoding/json"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
)

// list - Get all currencies
// @Summary Show all currencies
// @Description Get all currencies
// @Tags Currency
// @ID get-all-currencies
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} model.Currency
// @Router /currencies [get]
func list(w http.ResponseWriter, r *http.Request) {

	var currencies []model.Currency

	offset, limit := util.Paging(r.URL.Query())

	model.DB.Offset(offset).Limit(limit).Model(&model.Currency{}).Find(&currencies)

	json.NewEncoder(w).Encode(currencies)
}
