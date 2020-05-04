package item

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// details - Get cartItem by id
// @Summary Show a cartItem by id
// @Description Get cartItem by ID
// @Tags CartItem
// @ID get-cart-item-by-id
// @Produce  json
// @Param cart_id path string true "Cart ID"
// @Param item_id path string true "Cart-item ID"
// @Success 200 {object} model.CartItem
// @Failure 400 {array} string
// @Router /carts/{cart_id}/items/{item_id} [get]
func details(w http.ResponseWriter, r *http.Request) {
	cartItemID := chi.URLParam(r, "item_id")
	id, err := strconv.Atoi(cartItemID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	cartItem := &model.CartItem{}
	cartItem.ID = uint(id)

	err = model.DB.Model(&model.CartItem{}).First(&cartItem).Error
	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	model.DB.Model(&cartItem).Association("Product").Find(&cartItem.Product)
	model.DB.Model(&cartItem.Product).Association("Status").Find(&cartItem.Product.Status)
	model.DB.Model(&cartItem.Product).Association("ProductType").Find(&cartItem.Product.ProductType)
	model.DB.Model(&cartItem.Product).Association("Currency").Find(&cartItem.Product.Currency)

	util.Render(w, http.StatusOK, cartItem)
}
