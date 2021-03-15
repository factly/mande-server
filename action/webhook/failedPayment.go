package webhook

import (
	"encoding/json"
	"net/http"

	"github.com/factly/mande-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
)

// failedPayment - receive failed payment webhook
// @Summary receive failed payment webhook
// @Description receive failed payment webhook
// @Tags Webhook
// @ID failed-payment
// @Consume json
// @Produce  json
// @Param FailedPaymentEvent body model.FailedPaymentEvent true "Failed Payment object"
// @Success 201 {object} model.FailedPaymentEvent
// @Failure 400 {array} string
// @Router /webhooks/failed-payment [post]
func failedPayment(w http.ResponseWriter, r *http.Request) {
	// Verify the signature in webhook

	paymentEvent := model.FailedPaymentEvent{}

	err := json.NewDecoder(r.Body).Decode(&paymentEvent)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	razorpayOID := paymentEvent.Payload.Payment.Entity.OrderID
	razorpayPID := paymentEvent.Payload.Payment.Entity.ID

	tx := model.DB.Begin()

	var order model.Order
	var orderFound bool = true
	err = tx.Model(&model.Order{}).Where(&model.Order{
		RazorpayOrderID: razorpayOID,
	}).First(&order).Error

	if err != nil {
		orderFound = false
	}

	if !orderFound {
		err = tx.Model(&model.Membership{}).Where(&model.Membership{
			RazorpayOrderID: razorpayOID,
		}).UpdateColumn("status", "failed").Error
	} else {
		err = tx.Model(&order).UpdateColumn("status", "failed").Error
	}

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	err = tx.Model(&model.Payment{}).Where(&model.Payment{
		RazorpayPaymentID: razorpayPID,
	}).UpdateColumn("status", "failed").Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	renderx.JSON(w, http.StatusOK, nil)
}
