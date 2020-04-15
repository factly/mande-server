package actions

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/models"
	"github.com/go-chi/chi"
)

func GetPlan(w http.ResponseWriter, r *http.Request) {
	planId := chi.URLParam(r, "planId")
	id, err := strconv.Atoi(planId)

	if err != nil {
		log.Fatal(err)
	}

	req := &models.Plan{
		ID: uint(id),
	}
	json.NewDecoder(r.Body).Decode(&req)

	models.DB.Model(&models.Plan{}).First(&req)

	json.NewEncoder(w).Encode(req)
}

func CreatePlan(w http.ResponseWriter, r *http.Request) {

	req := &models.Plan{}

	json.NewDecoder(r.Body).Decode(&req)

	if validErrs := req.Validate(); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		w.Header().Set("Content-type", "applciation/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	err := models.DB.Model(&models.Plan{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(req)
}

func UpdatePlan(w http.ResponseWriter, r *http.Request) {

	planId := chi.URLParam(r, "planId")
	id, err := strconv.Atoi(planId)

	if err != nil {
		log.Fatal(err)
	}

	req := &models.Plan{
		ID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&req)
	plan := &models.Plan{}
	models.DB.First(&models.Plan{})

	if req.PlanName != "" {
		plan.PlanName = req.PlanName
	}
	if req.PlanInfo != "" {
		plan.PlanInfo = req.PlanInfo
	}
	if req.Status != "" {
		plan.Status = req.Status
	}

	models.DB.Model(&models.Plan{}).Update(&plan)

	json.NewEncoder(w).Encode(req)
}

func DeletePlan(w http.ResponseWriter, r *http.Request) {
	planId := chi.URLParam(r, "planId")
	id, err := strconv.Atoi(planId)

	if err != nil {
		log.Fatal(err)
	}

	plan := &models.Plan{
		ID: uint(id),
	}
	json.NewDecoder(r.Body).Decode(&plan)

	models.DB.First(&plan)
	models.DB.Delete(&plan)

	json.NewEncoder(w).Encode(plan)
}
