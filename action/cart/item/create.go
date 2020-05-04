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

// create - create cartItem
// @Summary Create cartItem
// @Description create cartItem
// @Tags CartItem
// @ID add-cart-item
// @Consume json
// @Produce  json
// @Param cart_id path string true "Cart ID"
// @Param CartItem body cartItem true "CartItem object"
// @Success 201 {object} model.CartItem
// @Failure 400 {array} string
// @Router /carts/{cart_id}/items [post]
func create(w http.ResponseWriter, r *http.Request) {

	cartItem := &model.CartItem{}

	json.NewDecoder(r.Body).Decode(&cartItem)

	validate := validator.New()
	err := validate.StructExcept(cartItem, "Product")
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}
	err = model.DB.Model(&model.CartItem{}).Create(&cartItem).Error

	if err != nil {
		log.Fatal(err)
	}
	model.DB.Model(&cartItem).Association("Product").Find(&cartItem.Product)
	model.DB.Model(&cartItem.Product).Association("Status").Find(&cartItem.Product.Status)
	model.DB.Model(&cartItem.Product).Association("ProductType").Find(&cartItem.Product.ProductType)
	model.DB.Model(&cartItem.Product).Association("Currency").Find(&cartItem.Product.Currency)

	util.Render(w, http.StatusCreated, cartItem)
}
