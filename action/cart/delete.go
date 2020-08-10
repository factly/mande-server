package cart

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// delete - Delete cart by id
// @Summary Delete a cart
// @Description Delete cart by ID
// @Tags Cart
// @ID delete-cart-by-id
// @Consume  json
// @Param cart_id path string true "Cart ID"
// @Success 200
// @Failure 400 {array} string
// @Router /carts/{cart_id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {

	cartID := chi.URLParam(r, "cart_id")
	id, err := strconv.Atoi(cartID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &model.Cart{}
	result.ID = uint(id)

	// check record exists or not
	err = model.DB.First(&result).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}
	model.DB.Delete(&result)

	renderx.JSON(w, http.StatusOK, nil)
}
