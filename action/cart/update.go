package cart

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/data-portal-server/util/meili"
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
// @Param X-User header string true "User ID"
// @Param cartitem_id path string true "Cart Item ID"
// @Param CartItem body cartitem false "Cart Item object"
// @Success 200 {object} model.CartItem
// @Failure 400 {array} string
// @Router /cartitems/{cartitem_id} [put]
func update(w http.ResponseWriter, r *http.Request) {
	cartitemID := chi.URLParam(r, "cartitem_id")
	id, err := strconv.Atoi(cartitemID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	uID, err := util.GetUser(r)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	cartitem := &cartitem{}

	err = json.NewDecoder(r.Body).Decode(&cartitem)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	validationError := validationx.Check(cartitem)
	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	result := &model.CartItem{}
	result.ID = uint(id)

	// check record exist or not
	err = model.DB.Model(&model.CartItem{}).First(&result).Error
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	tx := model.DB.Begin()

	err = tx.Model(&result).Updates(model.CartItem{
		Status:    cartitem.Status,
		UserID:    uint(uID),
		ProductID: cartitem.ProductID,
	}).Preload("Product").Preload("Product.Currency").Preload("Product.FeaturedMedium").Preload("Product.Tags").Preload("Product.Datasets").First(&result).Error

	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	// Update into meili index
	meiliObj := map[string]interface{}{
		"id":         result.ID,
		"kind":       "cartitem",
		"user_id":    result.UserID,
		"status":     result.Status,
		"product_id": result.ProductID,
	}

	err = meili.UpdateDocument(meiliObj)
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	tx.Commit()
	renderx.JSON(w, http.StatusOK, result)
}
