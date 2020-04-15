package actions

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/models"
	"github.com/go-chi/chi"
)

func GetMembership(w http.ResponseWriter, r *http.Request) {
	membershipId := chi.URLParam(r, "membershipId")
	id, err := strconv.Atoi(membershipId)
	if err != nil {
		log.Fatal(err)
	}

	req := &models.Membership{
		ID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&req)

	models.DB.Model(&models.Membership{}).First(&req)

	models.DB.Model(&req).Association("User").Find(&req.User)
	models.DB.Model(&req).Association("Plan").Find(&req.Plan)
	models.DB.Model(&req).Association("Payment").Find(&req.Payment)
	models.DB.Model(&req.Payment).Association("Currency").Find(&req.Payment.Currency)
	json.NewEncoder(w).Encode(req)
}

func CreateMembership(w http.ResponseWriter, r *http.Request) {

	req := &models.Membership{}
	json.NewDecoder(r.Body).Decode(&req)

	if validErrs := req.Validate(); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		w.Header().Set("Content-type", "applciation/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	err := models.DB.Model(&models.Membership{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}

	models.DB.Model(&req).Association("User").Find(&req.User)
	models.DB.Model(&req).Association("Plan").Find(&req.Plan)
	models.DB.Model(&req).Association("Payment").Find(&req.Payment)
	models.DB.Model(&req.Payment).Association("Currency").Find(&req.Payment.Currency)

	json.NewEncoder(w).Encode(req)
}

func UpdateMembership(w http.ResponseWriter, r *http.Request) {
	membershipId := chi.URLParam(r, "membershipId")
	id, err := strconv.Atoi(membershipId)
	if err != nil {
		log.Fatal(err)
	}

	req := &models.Membership{
		ID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&req)
	membership := &models.Membership{}
	models.DB.First(&req)
	membership = req

	if req.UserID != 0 {
		membership.UserID = req.UserID
	}
	if req.PaymentID != 0 {
		membership.PaymentID = req.PaymentID
	}
	if req.PlanID != 0 {
		membership.PlanID = req.PlanID
	}
	if req.Status != "" {
		membership.Status = req.Status
	}

	models.DB.Model(&models.Membership{}).Update(&membership)

	models.DB.Model(&membership).Association("User").Find(&membership.User)
	models.DB.Model(&membership).Association("Plan").Find(&membership.Plan)
	models.DB.Model(&membership).Association("Payment").Find(&membership.Payment)
	models.DB.Model(&membership.Payment).Association("Currency").Find(&membership.Payment.Currency)
	json.NewEncoder(w).Encode(membership)
}

func DeleteMembership(w http.ResponseWriter, r *http.Request) {
	membershipId := chi.URLParam(r, "membershipId")
	id, err := strconv.Atoi(membershipId)
	if err != nil {
		log.Fatal(err)
	}

	membership := &models.Membership{
		ID: uint(id),
	}
	json.NewDecoder(r.Body).Decode(&membership)

	models.DB.First(&membership)
	models.DB.Delete(&membership)

	json.NewEncoder(w).Encode(membership)
}
