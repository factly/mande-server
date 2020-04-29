package item

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
)

// detail - Get cartItem by id
// @Summary Show a cartItem by id
// @Description Get cartItem by ID
// @Tags CartItem
// @ID get-cartItem-by-id
// @Produce  json
// @Param cart_id path string true "Cart ID"
// @Param item_id path string true "Cart-item ID"
// @Param id path string true "CartItem ID"
// @Success 200 {object} model.CartItem
// @Failure 400 {array} string
// @Router /carts/cart_id/cart-items/{item_id} [get]
func detail(w http.ResponseWriter, r *http.Request) {

	cartItemID := chi.URLParam(r, "item_id")
	id, err := strconv.Atoi(cartItemID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &model.CartItem{
		ID: uint(id),
	}

	err = model.DB.Model(&model.CartItem{}).First(&req).Error
	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	model.DB.Model(&req).Association("Product").Find(&req.Product)
	model.DB.Model(&req.Product).Association("Status").Find(&req.Product.Status)
	model.DB.Model(&req.Product).Association("ProductType").Find(&req.Product.ProductType)
	model.DB.Model(&req.Product).Association("Currency").Find(&req.Product.Currency)
	json.NewEncoder(w).Encode(req)
}
