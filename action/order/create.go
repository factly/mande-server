package order

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/factly/x/renderx"
	"github.com/go-playground/validator/v10"
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

	validate := validator.New()
	err := validate.StructExcept(order, "Payment", "Cart")
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	result := &model.Order{
		UserID:    order.UserID,
		Status:    order.Status,
		PaymentID: order.PaymentID,
		CartID:    order.CartID,
	}

	err = model.DB.Model(&model.Order{}).Create(&result).Error

	if err != nil {
		log.Fatal(err)
	}
	model.DB.Model(&model.Order{}).Preload("Payment").Preload("Payment.Currency").Preload("Cart").First(&result)

	renderx.JSON(w, http.StatusCreated, result)
}
