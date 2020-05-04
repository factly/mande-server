package cart

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// details - Get cart by id
// @Summary Show a cart by id
// @Description Get cart by ID
// @Tags Cart
// @ID get-cart-by-id
// @Produce  json
// @Param cart_id path string true "Cart ID"
// @Success 200 {object} model.Cart
// @Failure 400 {array} string
// @Router /carts/{cart_id} [get]
func details(w http.ResponseWriter, r *http.Request) {

	cartID := chi.URLParam(r, "cart_id")
	id, err := strconv.Atoi(cartID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	cart := &model.Cart{}
	cart.ID = uint(id)

	err = model.DB.Model(&model.Cart{}).First(&cart).Error
	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	render.JSON(w, http.StatusOK, cart)
}
