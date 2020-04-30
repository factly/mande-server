package cart

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
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

	req := &model.Cart{
		ID: uint(id),
	}

	err = model.DB.Model(&model.Cart{}).First(&req).Error
	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(req)
}
