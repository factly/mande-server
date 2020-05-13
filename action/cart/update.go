package cart

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// update - Update cart by id
// @Summary Update a cart by id
// @Description Update cart by ID
// @Tags Cart
// @ID update-cart-by-id
// @Produce json
// @Consume json
// @Param cart_id path string true "Cart ID"
// @Param Cart body cart false "Cart"
// @Success 200 {object} model.Cart
// @Failure 400 {array} string
// @Router /carts/{cart_id} [put]
func update(w http.ResponseWriter, r *http.Request) {
	cartID := chi.URLParam(r, "cart_id")
	id, err := strconv.Atoi(cartID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	cart := &cart{}

	json.NewDecoder(r.Body).Decode(&cart)

	result := &model.Cart{}
	result.ID = uint(id)

	model.DB.Model(&result).Updates(model.Cart{
		Status: cart.Status,
		UserID: cart.UserID,
	}).First(&result)

	render.JSON(w, http.StatusOK, result)
}
