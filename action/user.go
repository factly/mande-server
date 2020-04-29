package action

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

// user request body
type user struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// GetUsers - Get all users
// @Summary Show all users
// @Description Get all users
// @Tags User
// @ID get-all-users
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} model.User
// @Router /users [get]
func GetUsers(w http.ResponseWriter, r *http.Request) {

	var users []model.User
	p := r.URL.Query().Get("page")
	pg, _ := strconv.Atoi(p) // pg contains page number
	l := r.URL.Query().Get("limit")
	li, _ := strconv.Atoi(l) // li contains perPage number

	offset := 0 // no. of records to skip
	limit := 5  // limt

	if li > 0 && li <= 10 {
		limit = li
	}

	if pg > 1 {
		offset = (pg - 1) * limit
	}

	model.DB.Offset(offset).Limit(limit).Model(&model.User{}).Find(&users)

	json.NewEncoder(w).Encode(users)
}

// GetUser - Get user by id
// @Summary Show a user by id
// @Description Get user by ID
// @Tags User
// @ID get-user-by-id
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} model.User
// @Failure 400 {array} string
// @Router /users/{id} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {

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

// CreateUser - Create user
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
func CreateUser(w http.ResponseWriter, r *http.Request) {

	req := &model.User{}
	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}
	err = model.DB.Model(&model.User{}).Create(&req).Error

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
// @Success 200 {object} model.User
// @Failure 400 {array} string
// @Router /users/{id} [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {

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

// DeleteUser - Delete user by id
// @Summary Delete a user
// @Description Delete user by ID
// @Tags User
// @ID delete-user-by-id
// @Consume  json
// @Param id path string true "User ID"
// @Success 200 {object} model.User
// @Failure 400 {array} string
// @Router /users/{id} [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {

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
