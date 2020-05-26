package item

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// update - Update cartItem by id
// @Summary Update a cartItem by id
// @Description Update cartItem by ID
// @Tags CartItem
// @ID update-cart-item-by-id
// @Produce json
// @Consume json
// @Param cart_id path string true "Cart ID"
// @Param item_id path string true "Cart-item ID"
// @Param CartItem body cartItem false "CartItem"
// @Success 200 {object} model.CartItem
// @Failure 400 {array} string
// @Router /carts/{cart_id}/items/{item_id} [put]
func update(w http.ResponseWriter, r *http.Request) {

	cartID := chi.URLParam(r, "cart_id")
	cid, _ := strconv.Atoi(cartID)

	cartItemID := chi.URLParam(r, "item_id")
	id, err := strconv.Atoi(cartItemID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	cartItem := &cartItem{}
	json.NewDecoder(r.Body).Decode(&cartItem)

	result := &model.CartItem{}
	result.ID = uint(id)
	result.CartID = uint(cid)

	model.DB.Model(&result).Updates(model.CartItem{
		ProductID: cartItem.ProductID,
	}).Preload("Product").Preload("Product.Status").Preload("Product.ProductType").Preload("Product.Currency").First(&result)

	render.JSON(w, http.StatusOK, result)
}
