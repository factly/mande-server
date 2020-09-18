package order

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

// create - create orders
// @Summary Create orders
// @Description create orders
// @Tags Order
// @ID add-orders
// @Consume json
// @Produce  json
// @Param Order body order true "Order object"
// @Success 201 {object} model.Order
// @Failure 400 {array} string
// @Router /orders [post]
func create(w http.ResponseWriter, r *http.Request) {

	order := &order{}

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	validationError := validationx.Check(order)
	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	result := &model.Order{
		UserID:    order.UserID,
		Status:    order.Status,
		PaymentID: order.PaymentID,
	}

	tx := model.DB.Begin()
	err = tx.Model(&model.Order{}).Create(&result).Error

	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	tx.Model(&model.Order{}).Preload("Payment").Preload("Payment.Currency").Preload("Cart").First(&result)

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
