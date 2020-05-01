package order

import (
	"encoding/json"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
)

// list - Get all orders
// @Summary Show all orders
// @Description Get all orders
// @Tags Order
// @ID get-all-orders
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} model.Order
// @Router /orders [get]
func list(w http.ResponseWriter, r *http.Request) {

	var orders []model.Order

	offset, limit := util.Paging(r.URL.Query())

	model.DB.Offset(offset).Limit(limit).Preload("Payment").Preload("Payment.Currency").Preload("Cart").Model(&model.Order{}).Find(&orders)

	json.NewEncoder(w).Encode(orders)
}
