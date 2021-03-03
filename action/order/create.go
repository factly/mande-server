package order

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/factly/mande-server/util/razorpay"

	"github.com/factly/mande-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/meilisearchx"
	"github.com/factly/x/middlewarex"
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
// @Param X-Organisation header string true "Organisation ID"
// @Success 201 {object} model.Order
// @Failure 400 {array} string
// @Router /orders [post]
func create(w http.ResponseWriter, r *http.Request) {
	uID, err := middlewarex.GetUser(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &model.Order{
		UserID: uint(uID),
		Status: "created",
	}

	tx := model.DB.WithContext(context.WithValue(r.Context(), userContext, uID)).Begin()

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

	var orderPrice int = 0

	// Delete cart items and append it in order & calculate price
	for _, item := range cartitems {
		result.Products = append(result.Products, *item.Product)

		if item.MembershipID == nil {
			orderPrice += item.Product.Price
		}

		if err := tx.Delete(item).Error; err != nil {
			tx.Rollback()
			loggerx.Error(err)
			errorx.Render(w, errorx.Parser(errorx.DBError()))
			return
		}
	}

	// Create order in database
	err = tx.Model(&model.Order{}).Create(&result).Error
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	// Fetch Currency
	currency := model.Currency{}
	model.DB.Model(&model.Currency{}).First(&currency)

	// Create a razorpay order and get razorpay orderID
	razorpayRequest := map[string]interface{}{
		"amount":          orderPrice * 100,
		"currency":        strings.ToUpper(currency.IsoCode),
		"receipt":         fmt.Sprint(result.ID),
		"payment_capture": 1,
	}

	orderBody, err := razorpay.Client.Order.Create(razorpayRequest, nil)
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	// Change order status to initiated and add razorpay_id in order table
	if _, found := orderBody["id"]; !found {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	result.Status = "processing"
	result.RazorpayOrderID = orderBody["id"].(string)

	err = tx.Omit("Products").Model(&result).Updates(result).Error
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	tx.Model(&model.Order{}).Preload("Products").Preload("Products.Datasets").Preload("Products.Tags").First(&result)

	// Insert into meili index
	meiliObj := map[string]interface{}{
		"id":         result.ID,
		"kind":       "order",
		"user_id":    result.UserID,
		"status":     result.Status,
		"payment_id": result.PaymentID,
	}

	err = meilisearchx.AddDocument("data-portal", meiliObj)
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	tx.Commit()

	renderx.JSON(w, http.StatusCreated, result)
}
