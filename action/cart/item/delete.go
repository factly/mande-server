package item

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
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
	cid, err := strconv.Atoi(cartID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	cartItemID := chi.URLParam(r, "item_id")
	id, err := strconv.Atoi(cartItemID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	cart := model.Cart{}
	cart.ID = uint(cid)

	product := model.Product{}
	product.ID = uint(id)

	// check record exists or not
	err = model.DB.Model(&cart).Association("Products").Find(&product).Error
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	model.DB.Model(&cart).Association("Products").Delete(&product)

	renderx.JSON(w, http.StatusOK, nil)
}
