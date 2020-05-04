package user

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-playground/validator/v10"
)

// create - Create user
// @Summary Create user
// @Description Create user
// @Tags User
// @ID add-user
// @Consume json
// @Produce  json
// @Param User body user true "User object"
// @Success 200 {object} model.User
// @Failure 400 {array} string
// @Router /users [post]
func create(w http.ResponseWriter, r *http.Request) {

	user := &model.User{}
	json.NewDecoder(r.Body).Decode(&user)

	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}
	err = model.DB.Model(&model.User{}).Create(&user).Error

	if err != nil {
		log.Fatal(err)
	}

	util.Render(w, http.StatusOK, user)
}
