package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
)

// delete - Delete user by id
// @Summary Delete a user
// @Description Delete user by ID
// @Tags User
// @ID delete-user-by-id
// @Consume  json
// @Param id path string true "User ID"
// @Success 200 {object} model.User
// @Failure 400 {array} string
// @Router /users/{id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {

	userID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	user := &model.User{
		ID: uint(id),
	}

	// check record exists or not
	err = model.DB.First(&user).Error
	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	model.DB.Delete(&user)

	json.NewEncoder(w).Encode(user)
}
