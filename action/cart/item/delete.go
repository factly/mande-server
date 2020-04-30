package item

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// deleteCartItem - Delete cartItem by id
// @Summary Delete a cartItem
// @Description Delete cartItem by ID
// @Tags CartItem
// @ID delete-cart-item-by-id
// @Consume  json
// @Param cart_id path string true "Cart ID"
// @Param item_id path string true "Cart-item ID"
// @Success 200 {object} model.CartItem
// @Failure 400 {array} string
// @Router /carts/{cart_id}/items/{item_id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {

	cartItemID := chi.URLParam(r, "item_id")
	id, err := strconv.Atoi(cartItemID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	cartItem := &model.CartItem{
		ID: uint(id),
	}

	// check record exists or not
	err = model.DB.First(&cartItem).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	model.DB.Delete(&cartItem)
	model.DB.Model(&cartItem).Association("Product").Find(&cartItem.Product)
	model.DB.Model(&cartItem.Product).Association("Status").Find(&cartItem.Product.Status)
	model.DB.Model(&cartItem.Product).Association("ProductType").Find(&cartItem.Product.ProductType)
	model.DB.Model(&cartItem.Product).Association("Currency").Find(&cartItem.Product.Currency)
	json.NewEncoder(w).Encode(cartItem)
}
