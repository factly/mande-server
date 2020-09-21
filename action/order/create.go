package order

import (
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/data-portal-server/util/meili"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
)

// create - create orders
// @Summary Create orders
// @Description create orders
// @Tags Order
// @ID add-orders
// @Consume json
// @Produce  json
// @Param X-User header string true "User ID"
// @Success 201 {object} model.Order
// @Failure 400 {array} string
// @Router /orders [post]
func create(w http.ResponseWriter, r *http.Request) {
	uID, err := util.GetUser(r)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &model.Order{
		UserID: uint(uID),
		Status: "created",
	}

	tx := model.DB.Begin()

	cartitems := make([]model.CartItem, 0)

	// Fetch all the items in cart
	tx.Model(&model.CartItem{}).Where(&model.CartItem{
		UserID: uint(uID),
	}).Preload("Product").Find(&cartitems)

	if len(cartitems) == 0 {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.CannotSaveChanges()))
		return
	}

	// Delete cart items
	for _, item := range cartitems {
		result.Products = append(result.Products, *item.Product)

		if err := tx.Delete(item).Error; err != nil {
			tx.Rollback()
			loggerx.Error(err)
			errorx.Render(w, errorx.Parser(errorx.DBError()))
			return
		}
	}

	// Create order in database
	err = tx.Model(&model.Order{}).Set("gorm:association_autoupdate", false).Create(&result).Error
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	// Create a razorpay order and get razorpay orderID
	// Change order status to initiated and add razorpay_id in order table

	tx.Model(&model.Order{}).Preload("Payment").Preload("Payment.Currency").Preload("Products").Preload("Products.Datasets").Preload("Products.Tags").First(&result)

	// Insert into meili index
	meiliObj := map[string]interface{}{
		"id":         result.ID,
		"kind":       "order",
		"user_id":    result.UserID,
		"status":     result.Status,
		"payment_id": result.PaymentID,
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
