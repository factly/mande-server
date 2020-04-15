package actions

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"../models"
	"github.com/go-chi/chi"
)

func GetTag(w http.ResponseWriter, r *http.Request) {
	tagId := chi.URLParam(r, "tagId")
	id, err := strconv.Atoi(tagId)
	if err != nil {
		log.Fatal(err)
	}

	req := &models.Tag{
		ID: uint(id),
	}
	json.NewDecoder(r.Body).Decode(&req)

	models.DB.Model(&models.Tag{}).First(&req)

	json.NewEncoder(w).Encode(req)
}

func CreateTag(w http.ResponseWriter, r *http.Request) {

	req := &models.Tag{}

	json.NewDecoder(r.Body).Decode(&req)

	if validErrs := req.Validate(); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		w.Header().Set("Content-type", "applciation/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	err := models.DB.Model(&models.Tag{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(req)
}

func UpdateTag(w http.ResponseWriter, r *http.Request) {
	tagId := chi.URLParam(r, "tagId")
	id, err := strconv.Atoi(tagId)
	if err != nil {
		log.Fatal(err)
	}

	req := &models.Tag{
		ID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&req)
	tag := &models.Tag{}
	models.DB.First(&models.Tag{})

	if req.Title != "" {
		tag.Title = req.Title
	}
	if req.Slug != "" {
		tag.Slug = req.Slug
	}

	models.DB.Model(&models.Tag{}).Update(&tag)

	json.NewEncoder(w).Encode(req)
}

func DeleteTag(w http.ResponseWriter, r *http.Request) {
	tagId := chi.URLParam(r, "tagId")
	id, err := strconv.Atoi(tagId)
	if err != nil {
		log.Fatal(err)
	}

	tag := &models.Tag{
		ID: uint(id),
	}
	json.NewDecoder(r.Body).Decode(&tag)

	models.DB.First(&tag)
	models.DB.Delete(&tag)

	json.NewEncoder(w).Encode(tag)
}
