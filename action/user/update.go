package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// update - Update user by id
// @Summary Update a user by id
// @Description Update user by ID
// @Tags User
// @ID update-user-by-id
// @Produce json
// @Consume json
// @Param user_id path string true "User ID"
// @Param User body user false "User"
// @Success 200 {object} model.User
// @Failure 400 {array} string
// @Router /users/{user_id} [put]
func update(w http.ResponseWriter, r *http.Request) {

	userID := chi.URLParam(r, "user_id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	user := &user{}

	json.NewDecoder(r.Body).Decode(&user)
	result := &model.User{}

	result.ID = uint(id)

	model.DB.Model(&result).Updates(&model.User{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}).First(&result)

	renderx.JSON(w, http.StatusOK, result)
}
