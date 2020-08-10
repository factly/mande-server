package item

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// details - Get cartItem by id
// @Summary Show a cartItem by id
// @Description Get cartItem by ID
// @Tags CartItem
// @ID get-cart-item-by-id
// @Produce  json
// @Param cart_id path string true "Cart ID"
// @Param item_id path string true "Cart-item ID"
// @Success 200 {object} model.CartItem
// @Failure 400 {array} string
// @Router /carts/{cart_id}/items/{item_id} [get]
func details(w http.ResponseWriter, r *http.Request) {

	cartID := chi.URLParam(r, "cart_id")
	cid, err := strconv.Atoi(cartID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	cartItemID := chi.URLParam(r, "item_id")
	id, err := strconv.Atoi(cartItemID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &model.CartItem{}
	result.ID = uint(id)
	result.CartID = uint(cid)

	err = model.DB.Model(&model.CartItem{}).First(&result).Error
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	model.DB.Preload("Product").Preload("Product.Currency").First(&result)

	renderx.JSON(w, http.StatusOK, result)
}
