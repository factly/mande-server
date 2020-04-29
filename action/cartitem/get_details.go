package cartitem

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
)

// getCartItemByID - Get cartItem by id
// @Summary Show a cartItem by id
// @Description Get cartItem by ID
// @Tags CartItem
// @ID get-cartItem-by-id
// @Produce  json
// @Param id path string true "CartItem ID"
// @Success 200 {object} model.CartItem
// @Failure 400 {array} string
// @Router /cart-items/{id} [get]
func getCartItemByID(w http.ResponseWriter, r *http.Request) {

	cartItemID := chi.URLParam(r, "id")
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
