package cart

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
)

// updateCart - Update cart by id
// @Summary Update a cart by id
// @Description Update cart by ID
// @Tags Cart
// @ID update-cart-by-id
// @Produce json
// @Consume json
// @Param id path string true "Cart ID"
// @Param Cart body cart false "Cart"
// @Success 200 {object} model.Cart
// @Failure 400 {array} string
// @Router /carts/{id} [put]
func updateCart(w http.ResponseWriter, r *http.Request) {

	cartID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(cartID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &model.Cart{}
	cart := &model.Cart{
		ID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&req)

	model.DB.Model(&cart).Updates(model.Cart{
		Status: req.Status,
		UserID: req.UserID,
	})
	model.DB.First(&cart)

	json.NewEncoder(w).Encode(cart)
}
