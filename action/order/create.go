package order

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
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

	order := &model.Order{}

	json.NewDecoder(r.Body).Decode(&order)

	validate := validator.New()
	err := validate.StructExcept(order, "Payment", "Cart")
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = model.DB.Model(&model.Order{}).Create(&order).Error

	if err != nil {
		log.Fatal(err)
	}
	model.DB.Model(&model.Order{}).Preload("Payment").Preload("Payment.Currency").Preload("Cart").First(&order)

	render.JSON(w, http.StatusCreated, order)
}
