package payment

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/factly/data-portal-server/util/razorpay"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/meili"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
)

// create - Create payment
// @Summary Create payment
// @Description Create payment
// @Tags Payment
// @ID add-payment
// @Consume json
// @Produce  json
// @Param Payment body payment true "Payment object"
// @Success 201 {object} model.Payment
// @Failure 400 {array} string
// @Router /payments [post]
func create(w http.ResponseWriter, r *http.Request) {

	payment := &payment{}
	err := json.NewDecoder(r.Body).Decode(&payment)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	validationError := validationx.Check(payment)
	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	result := &model.Payment{
		Gateway:           payment.Gateway,
		CurrencyID:        payment.CurrencyID,
		Status:            payment.Status,
		RazorpayPaymentID: payment.RazorpayPaymentID,
		RazorpaySignature: payment.RazorpaySignature,
	}

	if payment.For != "order" && payment.For != "membership" {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.Message{
			Code:    http.StatusUnprocessableEntity,
			Message: `"for" should be either "order" or "membership"`,
		}))
		return
	}

	tx := model.DB.Begin()

	// Fetch razorpay_order_id from entity for which the payment is created
	var razorpayOrderID string

	order := model.Order{}
	membership := model.Membership{}

	if payment.For == "order" {
		order.ID = payment.EntityID
		if err = tx.Model(&model.Order{}).First(&order).Error; err != nil {
			tx.Rollback()
			loggerx.Error(err)
			errorx.Render(w, errorx.Parser(errorx.CannotSaveChanges()))
			return
		}
		razorpayOrderID = order.RazorpayOrderID
	} else if payment.For == "membership" {
		membership.ID = payment.EntityID
		if err = tx.Model(&model.Membership{}).First(&membership).Error; err != nil {
			tx.Rollback()
			loggerx.Error(err)
			errorx.Render(w, errorx.Parser(errorx.CannotSaveChanges()))
			return
		}
		razorpayOrderID = membership.RazorpayOrderID
	}

	// verify the payment signature
	if !razorpay.VerifySignature(razorpayOrderID, payment.RazorpayPaymentID, payment.RazorpaySignature) {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.Message{
			Code:    http.StatusUnprocessableEntity,
			Message: `payment signature invalid`,
		}))
		return
	}

	// Get payment amount from razorpay
	razorpayPayment, err := razorpay.Client.Payment.Fetch(payment.RazorpayPaymentID, nil, nil)
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	if _, found := razorpayPayment["amount"]; !found {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	result.Amount = int(razorpayPayment["amount"].(float64) / 100)

	err = tx.Model(&model.Payment{}).Create(&result).Error

	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	// Update order/membership and mark as completed
	if payment.For == "order" {
		if err = tx.Model(&order).Updates(&model.Order{
			Status:    "complete",
			PaymentID: result.ID,
		}).Error; err != nil {
			tx.Rollback()
			loggerx.Error(err)
			errorx.Render(w, errorx.Parser(errorx.DBError()))
			return
		}
	} else if payment.For == "membership" {
		if err = tx.Model(&membership).Updates(&model.Membership{
			Status:    "complete",
			PaymentID: result.ID,
		}).Error; err != nil {
			tx.Rollback()
			loggerx.Error(err)
			errorx.Render(w, errorx.Parser(errorx.DBError()))
			return
		}
	}

	tx.Model(&result).Preload("Currency").First(&result)

	// Insert into meili index
	meiliObj := map[string]interface{}{
		"id":          result.ID,
		"kind":        "payment",
		"amount":      result.Amount,
		"gateway":     result.Gateway,
		"currency_id": result.CurrencyID,
		"status":      result.Status,
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
