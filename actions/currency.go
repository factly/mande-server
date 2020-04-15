package actions

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"../models"
	"github.com/go-chi/chi"
)

func GetCurrency(w http.ResponseWriter, r *http.Request) {
	currencyId := chi.URLParam(r, "currencyId")
	id, err := strconv.Atoi(currencyId)
	if err != nil {
		log.Fatal(err)
	}

	req := &models.Currency{
		ID: uint(id),
	}
	json.NewDecoder(r.Body).Decode(&req)

	models.DB.Model(&models.Currency{}).First(&req)

	json.NewEncoder(w).Encode(req)
}

func CreateCurrency(w http.ResponseWriter, r *http.Request) {

	req := &models.Currency{}

	json.NewDecoder(r.Body).Decode(&req)

	if validErrs := req.Validate(); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		w.Header().Set("Content-type", "applciation/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	err := models.DB.Model(&models.Currency{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(req)
}

func UpdateCurrency(w http.ResponseWriter, r *http.Request) {

	currencyId := chi.URLParam(r, "currencyId")
	id, err := strconv.Atoi(currencyId)
	if err != nil {
		log.Fatal(err)
	}

	req := &models.Currency{
		ID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&req)
	currency := &models.Currency{}
	models.DB.First(&models.Currency{})

	if req.IsoCode != "" {
		currency.IsoCode = req.IsoCode
	}
	if req.Name != "" {
		currency.Name = req.Name
	}

	models.DB.Model(&models.Currency{}).Update(&currency)

	json.NewEncoder(w).Encode(req)
}

func DeleteCurrency(w http.ResponseWriter, r *http.Request) {
	currencyId := chi.URLParam(r, "currencyId")
	id, err := strconv.Atoi(currencyId)
	if err != nil {
		log.Fatal(err)
	}

	currency := &models.Currency{
		ID: uint(id),
	}
	json.NewDecoder(r.Body).Decode(&currency)

	models.DB.First(&currency)
	models.DB.Delete(&currency)

	json.NewEncoder(w).Encode(currency)
}
