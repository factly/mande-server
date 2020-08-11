package item

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/factly/x/loggerx"

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
// @Success 201 {object} model.Product
// @Failure 400 {array} string
// @Router /carts/{cart_id}/items [post]
func create(w http.ResponseWriter, r *http.Request) {

	cartID := chi.URLParam(r, "cart_id")
	id, err := strconv.Atoi(cartID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	cartItem := &cartItem{}

	json.NewDecoder(r.Body).Decode(&cartItem)

	validationError := validationx.Check(cartItem)
	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	cart := model.Cart{}
	cart.ID = uint(id)

	// check if cart exist or not
	err = model.DB.Model(&model.Cart{}).First(&cart).Error
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	result := model.Product{}
	result.ID = cartItem.ProductID

	err = model.DB.Model(&cart).Association("Products").Append(&result).Error
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	renderx.JSON(w, http.StatusCreated, result)
}
