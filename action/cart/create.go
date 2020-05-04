package cart

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/data-portal-server/validation"
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
// @Success 200 {object} model.Cart
// @Failure 400 {array} string
// @Router /carts [post]
func create(w http.ResponseWriter, r *http.Request) {

	cart := &model.Cart{}

	json.NewDecoder(r.Body).Decode(&cart)

	validate := validator.New()
	err := validate.Struct(cart)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = model.DB.Model(&model.Cart{}).Create(&cart).Error

	if err != nil {
		log.Fatal(err)
	}

	util.Render(w, http.StatusOK, cart)
}
