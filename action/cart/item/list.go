package item

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
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
	id, err := strconv.Atoi(cartID)

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := paging{}
	result.Nodes = make([]model.CartItem, 0)

	offset, limit := paginationx.Parse(r.URL.Query())

	model.DB.Preload("Product").Preload("Product.Currency").Model(&model.CartItem{}).Where(&model.CartItem{CartID: uint(id)}).Count(&result.Total).Offset(offset).Limit(limit).Find(&result.Nodes)

	renderx.JSON(w, http.StatusOK, result)
}
