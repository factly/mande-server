package item

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
)

// list - Get all cartItems
// @Summary Show all cartItems
// @Description Get all cartItems
// @Tags CartItem
// @ID get-all-cartItems
// @Produce  json
// @Param cart_id path string true "Cart ID"
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} model.CartItem
// @Router /carts/{cart_id}/cart-items [get]
func list(w http.ResponseWriter, r *http.Request) {

	var cartItems []model.CartItem
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

	model.DB.Offset(offset).Limit(limit).Preload("Product").Preload("Product.Status").Preload("Product.ProductType").Preload("Product.Currency").Model(&model.CartItem{}).Find(&cartItems)

	json.NewEncoder(w).Encode(cartItems)
}
