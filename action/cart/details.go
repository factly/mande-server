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

// details - Get cart by id
// @Summary Show a cart by id
// @Description Get cart by ID
// @Tags Cart
// @ID get-cart-by-id
// @Produce  json
// @Param X-User header string true "User ID"
// @Param cartitem_id path string true "Cart Item ID"
// @Success 200 {object} model.CartItem
// @Failure 400 {array} string
// @Router /cartitems/{cartitem_id} [get]
func details(w http.ResponseWriter, r *http.Request) {

	cartitemID := chi.URLParam(r, "cartitem_id")
	id, err := strconv.Atoi(cartitemID)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &model.CartItem{}
	result.ID = uint(id)

	err = model.DB.Model(&model.CartItem{}).Preload("Product").Preload("Product.Currency").Preload("Product.FeaturedMedium").Preload("Product.Tags").Preload("Product.Datasets").First(&result).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	renderx.JSON(w, http.StatusOK, result)
}
