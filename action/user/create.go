package user

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
)

// create - Create user
// @Summary Create user
// @Description Create user
// @Tags User
// @ID add-user
// @Consume json
// @Produce  json
// @Param User body user true "User object"
// @Success 201 {object} model.User
// @Failure 400 {array} string
// @Router /users [post]
func create(w http.ResponseWriter, r *http.Request) {

	user := &user{}
	json.NewDecoder(r.Body).Decode(&user)

	err := validation.Validator.Struct(user)
	if err != nil {
		validation.ValidatorErrors(w, r, err)
		return
	}

	result := &model.User{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	err = model.DB.Model(&model.User{}).Create(&result).Error

	if err != nil {
		log.Fatal(err)
	}

	render.JSON(w, http.StatusCreated, result)
}
