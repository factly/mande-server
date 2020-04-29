package cartitem

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
)

// deleteCartItem - Delete cartItem by id
// @Summary Delete a cartItem
// @Description Delete cartItem by ID
// @Tags CartItem
// @ID delete-cartItem-by-id
// @Consume  json
// @Param id path string true "CartItem ID"
// @Success 200 {object} model.CartItem
// @Failure 400 {array} string
// @Router /cart-items/{id} [delete]
func deleteCartItem(w http.ResponseWriter, r *http.Request) {

	cartItemID := chi.URLParam(r, "id")
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
