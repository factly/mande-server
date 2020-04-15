package actions

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"../models"
	"github.com/go-chi/chi"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	id, err := strconv.Atoi(userId)
	if err != nil {
		log.Fatal(err)
	}

	req := &models.User{
		ID: uint(id),
	}
	json.NewDecoder(r.Body).Decode(&req)

	models.DB.Model(&models.User{}).First(&req)

	json.NewEncoder(w).Encode(req)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	req := &models.User{}
	json.NewDecoder(r.Body).Decode(&req)

	if validErrs := req.Validate(); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		w.Header().Set("Content-type", "applciation/json")
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

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	id, err := strconv.Atoi(userId)
	if err != nil {
		log.Fatal(err)
	}

	req := &models.User{
		ID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&req)
	user := &models.User{}
	models.DB.First(&models.User{})

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Name != "" {
		user.Name = req.Name
	}

	models.DB.Model(&models.User{}).Update(&user)

	json.NewEncoder(w).Encode(req)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	id, err := strconv.Atoi(userId)
	if err != nil {
		log.Fatal(err)
	}

	user := &models.User{
		ID: uint(id),
	}
	json.NewDecoder(r.Body).Decode(&user)

	models.DB.First(&user)
	models.DB.Delete(&user)

	json.NewEncoder(w).Encode(user)
}
