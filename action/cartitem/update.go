package cartitem

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
)

// updateCartItem - Update cartItem by id
// @Summary Update a cartItem by id
// @Description Update cartItem by ID
// @Tags CartItem
// @ID update-cartItem-by-id
// @Produce json
// @Consume json
// @Param id path string true "CartItem ID"
// @Param CartItem body cartItem false "CartItem"
// @Success 200 {object} model.CartItem
// @Failure 400 {array} string
// @Router /cart-items/{id} [put]
func updateCartItem(w http.ResponseWriter, r *http.Request) {

	cartItemID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(cartItemID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &model.CartItem{}
	cartItem := &model.CartItem{
		ID: uint(id),
	}

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
