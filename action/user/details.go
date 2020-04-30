package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// details - Get user by id
// @Summary Show a user by id
// @Description Get user by ID
// @Tags User
// @ID get-user-by-id
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} model.User
// @Failure 400 {array} string
// @Router /users/{id} [get]
func details(w http.ResponseWriter, r *http.Request) {

	userID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &model.User{
		ID: uint(id),
	}

	err = model.DB.Model(&model.User{}).First(&req).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(req)
}
