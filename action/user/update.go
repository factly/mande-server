package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
)

// updateUser - Update user by id
// @Summary Update a user by id
// @Description Update user by ID
// @Tags User
// @ID update-user-by-id
// @Produce json
// @Consume json
// @Param id path string true "User ID"
// @Param User body user false "User"
// @Success 200 {object} model.User
// @Failure 400 {array} string
// @Router /users/{id} [put]
func updateUser(w http.ResponseWriter, r *http.Request) {

	userID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &model.User{}

	json.NewDecoder(r.Body).Decode(&req)
	user := &model.User{
		ID: uint(id),
	}

	model.DB.Model(&user).Updates(&model.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	model.DB.First(&user)
	json.NewEncoder(w).Encode(user)
}
