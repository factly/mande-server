package item

import (
	"encoding/json"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
)

// list - Get all cartItems
// @Summary Show all cartItems
// @Description Get all cartItems
// @Tags CartItem
// @ID get-all-cart-items
// @Produce  json
// @Param cart_id path string true "Cart ID"
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} model.CartItem
// @Router /carts/{cart_id}/items [get]
func list(w http.ResponseWriter, r *http.Request) {

	var cartItems []model.CartItem

	offset, limit := util.Paging(r.URL.Query())

	model.DB.Offset(offset).Limit(limit).Preload("Product").Preload("Product.Status").Preload("Product.ProductType").Preload("Product.Currency").Model(&model.CartItem{}).Find(&cartItems)

	json.NewEncoder(w).Encode(cartItems)
}
