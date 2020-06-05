package item

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/factly/x/renderx"
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
// @Success 200
// @Failure 400 {array} string
// @Router /carts/{cart_id}/items/{item_id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {

	cartID := chi.URLParam(r, "cart_id")
	cid, _ := strconv.Atoi(cartID)

	cartItemID := chi.URLParam(r, "item_id")
	id, err := strconv.Atoi(cartItemID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	result := &model.CartItem{}
	result.ID = uint(id)
	result.CartID = uint(cid)

	// check record exists or not
	err = model.DB.First(&result).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	model.DB.Delete(&result)

	renderx.JSON(w, http.StatusOK, nil)
}
