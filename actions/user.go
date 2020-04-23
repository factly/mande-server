package actions

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/models"
	"github.com/factly/data-portal-api/validationerrors"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

// user request body
type user struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// GetUser - Get user by id
// @Summary Show a user by id
// @Description Get user by ID
// @Tags User
// @ID get-user-by-id
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {array} string
// @Router /users/{id} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		validationerrors.InvalidID(w, r)
		return
	}

	req := &models.User{
		ID: uint(id),
	}

	models.DB.Model(&models.User{}).First(&req)

	json.NewEncoder(w).Encode(req)
}

// CreateUser - Create user
// @Summary Create user
// @Description Create user
// @Tags User
// @ID add-user
// @Consume json
// @Produce  json
// @Param User body user true "User object"
// @Success 200 {object} models.User
// @Failure 400 {array} string
// @Router /users [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	req := &models.User{}
	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		msg := err.Error()
		validationerrors.ValidErrors(w, r, msg)
		return
	}
	err = models.DB.Model(&models.User{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(req)
}

// UpdateUser - Update user by id
// @Summary Update a user by id
// @Description Update user by ID
// @Tags User
// @ID update-user-by-id
// @Produce json
// @Consume json
// @Param id path string true "User ID"
// @Param User body user false "User"
// @Success 200 {object} models.User
// @Failure 400 {array} string
// @Router /users/{id} [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		validationerrors.InvalidID(w, r)
		return
	}

	req := &models.User{}

	json.NewDecoder(r.Body).Decode(&req)
	user := &models.User{
		ID: uint(id),
	}

	models.DB.Model(&user).Updates(&models.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	models.DB.First(&user)
	json.NewEncoder(w).Encode(user)
}

// DeleteUser - Delete user by id
// @Summary Delete a user
// @Description Delete user by ID
// @Tags User
// @ID delete-user-by-id
// @Consume  json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {array} string
// @Router /users/{id} [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		validationerrors.InvalidID(w, r)
		return
	}

	user := &models.User{
		ID: uint(id),
	}

	// check record exists or not
	err = models.DB.First(&user).Error
	if err != nil {
		validationerrors.RecordNotFound(w, r)
		return
	}

	models.DB.Delete(&user)

	json.NewEncoder(w).Encode(user)
}
