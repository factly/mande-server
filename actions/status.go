package actions

import (
	"encoding/json"
	"log"
	"net/http"

	"../models"
)

func CreateStatus(w http.ResponseWriter, r *http.Request) {

	req := &models.Status{}
	json.NewDecoder(r.Body).Decode(&req)

	if validErrs := req.Validate(); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		w.Header().Set("Content-type", "applciation/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	err := models.DB.Model(&models.Status{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(req)
}

func UpdateStatus(w http.ResponseWriter, r *http.Request) {

	req := &models.Status{}

	json.NewDecoder(r.Body).Decode(&req)
	status := &models.Status{}
	models.DB.First(&models.Status{})

	if req.Name != "" {
		status.Name = req.Name
	}

	models.DB.Model(&models.Status{}).Update(&status)

	json.NewEncoder(w).Encode(req)
}

func DeleteStatus(w http.ResponseWriter, r *http.Request) {
	status := &models.Status{}
	json.NewDecoder(r.Body).Decode(&status)

	models.DB.First(&status)
	models.DB.Delete(&status)

	json.NewEncoder(w).Encode(status)
}
