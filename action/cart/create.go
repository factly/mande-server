package cart

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/meili"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/middlewarex"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
)

// create - create cart
// @Summary Create cart
// @Description create cart
// @Tags Cart
// @ID add-cart
// @Consume json
// @Produce  json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param CartItem body cartitem true "Cart Item object"
// @Success 201 {object} model.CartItem
// @Failure 400 {array} string
// @Router /cartitems [post]
func create(w http.ResponseWriter, r *http.Request) {
	uID, err := middlewarex.GetUser(r.Context())
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

	membershipID := &cartitem.MembershipID
	if cartitem.MembershipID == 0 {
		membershipID = nil
	}

	result := &model.CartItem{
		Status:       cartitem.Status,
		UserID:       uint(uID),
		ProductID:    cartitem.ProductID,
		MembershipID: membershipID,
	}

	tx := model.DB.WithContext(context.WithValue(r.Context(), userContext, uID)).Begin()

	if cartitem.MembershipID != 0 {
		membership := model.Membership{}
		membership.ID = cartitem.MembershipID
		err = tx.Model(&model.Membership{}).First(&membership).Error
		if err != nil {
			tx.Rollback()
			loggerx.Error(err)
			errorx.Render(w, errorx.Parser(errorx.CannotSaveChanges()))
			return
		}
	}

	err = tx.Model(&model.CartItem{}).Create(&result).Error
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	tx.Preload("Product").Preload("Membership").Preload("Membership.Plan").Preload("Product.Currency").Preload("Product.FeaturedMedium").Preload("Product.Tags").Preload("Product.Datasets").First(&result)

	// Insert into meili index
	meiliObj := map[string]interface{}{
		"id":         result.ID,
		"kind":       "cartitem",
		"user_id":    result.UserID,
		"status":     result.Status,
		"product_id": result.ProductID,
	}

	err = meili.AddDocument(meiliObj)
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	tx.Commit()

	renderx.JSON(w, http.StatusCreated, result)
}
