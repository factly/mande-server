package actions

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/models"
	"github.com/go-chi/chi"
)

func GetPayment(w http.ResponseWriter, r *http.Request) {
	paymentId := chi.URLParam(r, "paymentId")
	id, err := strconv.Atoi(paymentId)

	if err != nil {
		log.Fatal(err)
	}

	req := &models.Payment{
		ID: uint(id),
	}

	models.DB.Model(&models.Payment{}).First(&req)

	models.DB.Model(&req).Association("Currency").Find(&req.Currency)
	json.NewEncoder(w).Encode(req)
}

func CreatePayment(w http.ResponseWriter, r *http.Request) {

	req := &models.Payment{}
	json.NewDecoder(r.Body).Decode(&req)

	if validErrs := req.Validate(); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		w.Header().Set("Content-type", "applciation/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	err := models.DB.Model(&models.Payment{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}
	models.DB.Model(&req).Association("Currency").Find(&req.Currency)
	json.NewEncoder(w).Encode(req)
}

func UpdatePayment(w http.ResponseWriter, r *http.Request) {
	paymentId := chi.URLParam(r, "paymentId")
	id, err := strconv.Atoi(paymentId)

	if err != nil {
		log.Fatal(err)
	}

	req := &models.Payment{
		ID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&req)
	payment := &models.Payment{}
	models.DB.First(&models.Payment{})

	if req.Amount != 0 {
		payment.Amount = req.Amount
	}
	if req.Gateway != "" {
		payment.Gateway = req.Gateway
	}
	if req.Status != "" {
		payment.Status = req.Status
	}
	if req.CurrencyID != 0 {
		payment.CurrencyID = req.CurrencyID
	}

	models.DB.Model(&models.Payment{}).Update(&payment)

	json.NewEncoder(w).Encode(req)
}

func DeletePayment(w http.ResponseWriter, r *http.Request) {
	paymentId := chi.URLParam(r, "paymentId")
	id, err := strconv.Atoi(paymentId)

	if err != nil {
		log.Fatal(err)
	}

	payment := &models.Payment{
		ID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&payment)

	models.DB.First(&payment)
	models.DB.Delete(&payment)

	json.NewEncoder(w).Encode(payment)
}
