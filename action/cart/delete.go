package cart

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
)

// deleteCart - Delete cart by id
// @Summary Delete a cart
// @Description Delete cart by ID
// @Tags Cart
// @ID delete-cart-by-id
// @Consume  json
// @Param id path string true "Cart ID"
// @Success 200 {object} model.Cart
// @Failure 400 {array} string
// @Router /carts/{id} [delete]
func deleteCart(w http.ResponseWriter, r *http.Request) {

	cartID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(cartID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	cart := &model.Cart{
		ID: uint(id),
	}

	// check record exists or not
	err = model.DB.First(&cart).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	model.DB.Delete(&cart)

	json.NewEncoder(w).Encode(cart)
}
