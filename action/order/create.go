package order

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-playground/validator/v10"
)

// createOrder - create orders
// @Summary Create orders
// @Description create orders
// @Tags Order
// @ID add-orders
// @Consume json
// @Produce  json
// @Param Order body order true "Order object"
// @Success 200 {object} model.Order
// @Failure 400 {array} string
// @Router /orders [post]
func createOrder(w http.ResponseWriter, r *http.Request) {

	req := &model.Order{}

	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.StructExcept(req, "Payment", "Cart")
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = model.DB.Model(&model.Order{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}
	model.DB.Model(&model.Order{}).Preload("Payment").Preload("Payment.Currency").Preload("Cart").First(&req)
	json.NewEncoder(w).Encode(req)
}