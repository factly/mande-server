package user

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// details - Get user by id
// @Summary Show a user by id
// @Description Get user by ID
// @Tags User
// @ID get-user-by-id
// @Produce  json
// @Param user_id path string true "User ID"
// @Success 200 {object} model.User
// @Failure 400 {array} string
// @Router /users/{user_id} [get]
func details(w http.ResponseWriter, r *http.Request) {

	userID := chi.URLParam(r, "user_id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	user := &model.User{}

	user.ID = uint(id)

	err = model.DB.Model(&model.User{}).First(&user).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	render.JSON(w, http.StatusOK, user)
}
