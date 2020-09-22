package payment

import (
	"encoding/json"
	"errors"
	"net/http"

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
		Amount:     payment.Amount,
		Gateway:    payment.Gateway,
		CurrencyID: payment.CurrencyID,
		Status:     payment.Status,
	}

	// verify the payment signature

	tx := model.DB.Begin()
	err = tx.Model(&model.Payment{}).Create(&result).Error

	// Get the order or membership which has razorpay_order_id == payment.RazorpayOrderID
	// change order/membership status to complete and update it to add the created payment id

	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
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
