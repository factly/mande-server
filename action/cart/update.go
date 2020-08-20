package cart

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
	"github.com/go-chi/chi"
)

// update - Update cart by id
// @Summary Update a cart by id
// @Description Update cart by ID
// @Tags Cart
// @ID update-cart-by-id
// @Produce json
// @Consume json
// @Param cart_id path string true "Cart ID"
// @Param Cart body cart false "Cart"
// @Success 200 {object} model.Cart
// @Failure 400 {array} string
// @Router /carts/{cart_id} [put]
func update(w http.ResponseWriter, r *http.Request) {
	cartID := chi.URLParam(r, "cart_id")
	id, err := strconv.Atoi(cartID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	cart := &cart{}

	err = json.NewDecoder(r.Body).Decode(&cart)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	validationError := validationx.Check(cart)
	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	result := &model.Cart{}
	result.ID = uint(id)
	result.Products = make([]model.Product, 0)

	// check record exist or not
	err = model.DB.Model(&model.Cart{}).Preload("Products").First(&result).Error
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	oldProducts := result.Products
	newProducts := make([]model.Product, 0)
	model.DB.Model(&model.Product{}).Where(cart.ProductIDs).Find(&newProducts)

	if len(oldProducts) > 0 {
		model.DB.Model(&result).Association("Products").Delete(oldProducts)
	}
	if len(newProducts) == 0 {
		newProducts = nil
	}

	model.DB.Model(&result).Set("gorm:association_autoupdate", false).Updates(model.Cart{
		Status:   cart.Status,
		UserID:   cart.UserID,
		Products: newProducts,
	}).Preload("Products").Preload("Products.Tags").Preload("Products.Datasets").Preload("Products.Currency").First(&result)

	renderx.JSON(w, http.StatusOK, result)
}
