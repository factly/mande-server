package actions

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/models"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

// Plan request body
type plan struct {
	PlanName string `json:"plan_name"`
	PlanInfo string `json:"plan_info"`
	Status   string `json:"status"`
}

// GetPlan - Get plan by id
// @Summary Show a plan by id
// @Description Get plan by ID
// @Tags Plan
// @ID get-plan-by-id
// @Produce  json
// @Param id path string true "Plan ID"
// @Success 200 {object} models.Plan
// @Failure 400 {array} string
// @Router /plans/{id} [get]
func GetPlan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	planID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(planID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &models.Plan{
		ID: uint(id),
	}

	models.DB.Model(&models.Plan{}).First(&req)

	json.NewEncoder(w).Encode(req)
}

// CreatePlan - create plan
// @Summary Create plan
// @Description create plan
// @Tags Plan
// @ID add-plan
// @Consume json
// @Produce  json
// @Param Plan body plan true "Plan object"
// @Success 200 {object} models.Plan
// @Router /plans [post]
func CreatePlan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	req := &models.Plan{}

	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = models.DB.Model(&models.Plan{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(req)
}

// UpdatePlan - Update plan by id
// @Summary Update a plan by id
// @Description Update plan by ID
// @Tags Plan
// @ID update-plan-by-id
// @Produce json
// @Consume json
// @Param id path string true "Plan ID"
// @Param Plan body plan false "Plan"
// @Success 200 {object} models.Plan
// @Failure 400 {array} string
// @Router /plans/{id} [put]
func UpdatePlan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	planID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(planID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &models.Plan{}
	plan := &models.Plan{
		ID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&req)

	models.DB.Model(&plan).Updates(models.Plan{
		PlanName: req.PlanName,
		PlanInfo: req.PlanInfo,
		Status:   req.Status,
	})
	models.DB.First(&plan)

	json.NewEncoder(w).Encode(plan)
}

// DeletePlan - Delete plan by id
// @Summary Delete a plan
// @Description Delete plan by ID
// @Tags Plan
// @ID delete-plan-by-id
// @Consume  json
// @Param id path string true "Plan ID"
// @Success 200 {object} models.Plan
// @Failure 400 {array} string
// @Router /plans/{id} [delete]
func DeletePlan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	planID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(planID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	plan := &models.Plan{
		ID: uint(id),
	}

	// check record exists or not
	err = models.DB.First(&plan).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	models.DB.Delete(&plan)

	json.NewEncoder(w).Encode(plan)
}
