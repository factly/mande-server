package item

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
	"github.com/go-chi/chi"
)

// create - create cartItem
// @Summary Create cartItem
// @Description create cartItem
// @Tags CartItem
// @ID add-cart-item
// @Consume json
// @Produce  json
// @Param cart_id path string true "Cart ID"
// @Param CartItem body cartItem true "CartItem object"
// @Success 201 {object} model.CartItem
// @Failure 400 {array} string
// @Router /carts/{cart_id}/items [post]
func create(w http.ResponseWriter, r *http.Request) {

	cartID := chi.URLParam(r, "cart_id")
	id, err := strconv.Atoi(cartID)

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	cartItem := &cartItem{}
	result := &model.CartItem{}

	json.NewDecoder(r.Body).Decode(&cartItem)

	validationError := validationx.Check(cartItem)
	if validationError != nil {
		errorx.Render(w, validationError)
		return
	}

	result.CartID = uint(id)
	result.ProductID = cartItem.ProductID

	err = model.DB.Model(&model.CartItem{}).Create(&result).Error

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	model.DB.Preload("Product").Preload("Product.Currency").First(&result)

	renderx.JSON(w, http.StatusCreated, result)
}
