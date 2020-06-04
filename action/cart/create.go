package cart

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
)

// create - create cart
// @Summary Create cart
// @Description create cart
// @Tags Cart
// @ID add-cart
// @Consume json
// @Produce  json
// @Param Cart body cart true "Cart object"
// @Success 201 {object} model.Cart
// @Failure 400 {array} string
// @Router /carts [post]
func create(w http.ResponseWriter, r *http.Request) {

	cart := &cart{}

	json.NewDecoder(r.Body).Decode(&cart)

	err := validation.Validator.Struct(cart)
	if err != nil {
		validation.ValidatorErrors(w, r, err)
		return
	}

	result := &model.Cart{
		Status: cart.Status,
		UserID: cart.UserID,
	}

	err = model.DB.Model(&model.Cart{}).Create(&result).Error

	if err != nil {
		log.Fatal(err)
	}

	render.JSON(w, http.StatusCreated, result)
}
