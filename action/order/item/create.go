package item

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-playground/validator/v10"
)

// create - create order items
// @Summary Create order items
// @Description create order items
// @Tags OrderItem
// @ID add-order-item
// @Consume json
// @Produce  json
// @Param order_id path string true "Order ID"
// @Param OrderItem body orderItem true "Order item object"
// @Success 200 {object} model.OrderItem
// @Router /orders/{order_id}/items [post]
func create(w http.ResponseWriter, r *http.Request) {

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

	util.Render(w, http.StatusOK, req)
}
