package cart

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/factly/x/renderx"
	"github.com/go-playground/validator/v10"
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

	validate := validator.New()
	err := validate.Struct(cart)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
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

	renderx.JSON(w, http.StatusCreated, result)
}
