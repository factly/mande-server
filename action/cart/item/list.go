package item

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/data-portal-server/util/render"
	"github.com/go-chi/chi"
)

// list response
type paging struct {
	Total int              `json:"total"`
	Nodes []model.CartItem `json:"nodes"`
}

// list - Get all cartItems
// @Summary Show all cartItems
// @Description Get all cartItems
// @Tags CartItem
// @ID get-all-cart-items
// @Produce  json
// @Param cart_id path string true "Cart ID"
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /carts/{cart_id}/items [get]
func list(w http.ResponseWriter, r *http.Request) {

	cartID := chi.URLParam(r, "cart_id")
	id, _ := strconv.Atoi(cartID)

	result := paging{}

	offset, limit := util.Paging(r.URL.Query())

	model.DB.Preload("Product").Preload("Product.ProductType").Preload("Product.Currency").Model(&model.CartItem{}).Where(&model.CartItem{CartID: uint(id)}).Count(&result.Total).Offset(offset).Limit(limit).Find(&result.Nodes)

	render.JSON(w, http.StatusOK, result)
}
