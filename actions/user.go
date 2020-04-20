package actions

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/models"
	"github.com/go-chi/chi"
)

// user request body
type user struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

// GetUser - Get user by id
// @Summary Show a user by id
// @Description Get user by ID
// @Tags User
// @ID get-user-by-id
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Router /users/{id} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		log.Fatal(err)
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
// @Router /users [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	req := &models.User{}
	json.NewDecoder(r.Body).Decode(&req)

	if validErrs := req.Validate(); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	err := models.DB.Model(&models.User{}).Create(&req).Error

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
// @Router /users/{id} [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		log.Fatal(err)
	}

	req := &models.User{}

	json.NewDecoder(r.Body).Decode(&req)
	user := &models.User{
		ID: uint(id),
	}

	models.DB.Model(&user).Updates(&models.User{Email: req.Email, Name: req.Name})
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
// @Router /users/{id} [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(userID)
	if err != nil {
		log.Fatal(err)
	}

	user := &models.User{
		ID: uint(id),
	}

	models.DB.First(&user)
	models.DB.Delete(&user)

	json.NewEncoder(w).Encode(user)
}
