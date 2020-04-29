package order

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
)

// getOrders - Get all orders
// @Summary Show all orders
// @Description Get all orders
// @Tags Order
// @ID get-all-orders
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} model.Order
// @Router /orders [get]
func getOrders(w http.ResponseWriter, r *http.Request) {

	var orders []model.Order
	p := r.URL.Query().Get("page")
	pg, _ := strconv.Atoi(p) // pg contains page number
	l := r.URL.Query().Get("limit")
	li, _ := strconv.Atoi(l) // li contains perPage number

	offset := 0 // no. of records to skip
	limit := 5  // limt

	if li > 0 && li <= 10 {
		limit = li
	}

	if pg > 1 {
		offset = (pg - 1) * limit
	}

	model.DB.Offset(offset).Limit(limit).Preload("Payment").Preload("Payment.Currency").Preload("Cart").Model(&model.Order{}).Find(&orders)

	json.NewEncoder(w).Encode(orders)
}
