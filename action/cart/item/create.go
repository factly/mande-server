package item

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
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

	cartID := chi.URLParam(r, "cart_id")
	id, _ := strconv.Atoi(cartID)

	cartItem := &cartItem{}
	result := &model.CartItem{}

	json.NewDecoder(r.Body).Decode(&cartItem)

	validate := validator.New()
	err := validate.StructExcept(cartItem, "Product")

	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	result.CartID = uint(id)
	result.ProductID = cartItem.ProductID

	err = model.DB.Model(&model.CartItem{}).Create(&result).Error

	if err != nil {
		log.Fatal(err)
	}

	model.DB.Preload("Product").Preload("Product.ProductType").Preload("Product.Currency").First(&result)

	renderx.JSON(w, http.StatusCreated, result)
}
