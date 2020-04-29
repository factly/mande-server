package item

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-playground/validator/v10"
)

// createOrderItem - create order items
// @Summary Create order items
// @Description create order items
// @Tags OrderItem
// @ID add-order-item
// @Consume json
// @Produce  json
// @Param id path string true "Order ID"
// @Param OrderItem body orderItem true "Order item object"
// @Success 200 {object} model.OrderItem
// @Router /orders/{id}/order-items [post]
func createOrderItem(w http.ResponseWriter, r *http.Request) {

	req := &model.OrderItem{}

	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.StructExcept(req, "Product", "Order")
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = model.DB.Model(&model.OrderItem{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}
	model.DB.Model(&model.OrderItem{}).Preload("Product").Preload("Product.Status").Preload("Product.ProductType").Preload("Product.Currency").Preload("Order").First(&req)
	json.NewEncoder(w).Encode(req)
}