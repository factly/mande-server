package item

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
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

	cartItemID := chi.URLParam(r, "item_id")
	id, err := strconv.Atoi(cartItemID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &model.CartItem{}
	cartItem := &model.CartItem{}
	cartItem.ID = uint(id)

	json.NewDecoder(r.Body).Decode(&req)

	model.DB.Model(&cartItem).Updates(model.CartItem{
		IsDeleted: req.IsDeleted,
		CartID:    req.CartID,
		ProductID: req.ProductID,
	})
	model.DB.First(&cartItem)
	model.DB.Model(&cartItem).Association("Product").Find(&cartItem.Product)
	model.DB.Model(&cartItem.Product).Association("Status").Find(&cartItem.Product.Status)
	model.DB.Model(&cartItem.Product).Association("ProductType").Find(&cartItem.Product.ProductType)
	model.DB.Model(&cartItem.Product).Association("Currency").Find(&cartItem.Product.Currency)
	json.NewEncoder(w).Encode(cartItem)
}
