package item

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-playground/validator/v10"
)

// create - create cartItem
// @Summary Create cartItem
// @Description create cartItem
// @Tags CartItem
// @ID add-cart-item
// @Consume json
// @Produce  json
// @Param cart_id path string true "Cart ID"
// @Param CartItem body cartItem true "CartItem object"
// @Success 200 {object} model.CartItem
// @Failure 400 {array} string
// @Router /carts/{cart_id}/cart-items [post]
func create(w http.ResponseWriter, r *http.Request) {

	req := &model.CartItem{}

	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.StructExcept(req, "Product")
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}
	err = model.DB.Model(&model.CartItem{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}
	model.DB.Model(&req).Association("Product").Find(&req.Product)
	model.DB.Model(&req.Product).Association("Status").Find(&req.Product.Status)
	model.DB.Model(&req.Product).Association("ProductType").Find(&req.Product.ProductType)
	model.DB.Model(&req.Product).Association("Currency").Find(&req.Product.Currency)
	json.NewEncoder(w).Encode(req)
}
