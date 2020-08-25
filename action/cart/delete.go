package cart

import (
	"errors"
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
	err = model.DB.Preload("Products").First(&result).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	// check if cart is associated with order
	var totAssociated int
	model.DB.Model(&model.Order{}).Where(&model.Order{
		CartID: uint(id),
	}).Count(&totAssociated)

	if totAssociated != 0 {
		loggerx.Error(errors.New("cart is associated with order"))
		errorx.Render(w, errorx.Parser(errorx.CannotSaveChanges()))
		return
	}

	tx := model.DB.Begin()

	// delete associations
	if len(result.Products) > 0 {
		err = tx.Model(&result).Association("Products").Delete(&result.Products).Error
		if err != nil {
			tx.Rollback()
			loggerx.Error(err)
			errorx.Render(w, errorx.Parser(errorx.DBError()))
			return
		}
	}

	err = tx.Delete(&result).Error
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}
	tx.Commit()

	renderx.JSON(w, http.StatusOK, nil)
}
