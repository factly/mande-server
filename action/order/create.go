package order

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
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

	json.NewDecoder(r.Body).Decode(&order)

	validationError := validationx.Check(order)
	if validationError != nil {
		renderx.JSON(w, http.StatusBadRequest, validationError)
		return
	}

	result := &model.Order{
		UserID:    order.UserID,
		Status:    order.Status,
		PaymentID: order.PaymentID,
		CartID:    order.CartID,
	}

	err := model.DB.Model(&model.Order{}).Create(&result).Error

	if err != nil {
		log.Fatal(err)
	}
	model.DB.Model(&model.Order{}).Preload("Payment").Preload("Payment.Currency").Preload("Cart").First(&result)

	renderx.JSON(w, http.StatusCreated, result)
}
