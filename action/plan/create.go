package plan

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
)

// Create - create plan
// @Summary Create plan
// @Description create plan
// @Tags Plan
// @ID add-plan
// @Consume json
// @Produce  json
// @Param Plan body plan true "Plan object"
// @Success 201 {object} model.Plan
// @Router /plans [post]
func Create(w http.ResponseWriter, r *http.Request) {

	plan := &plan{}

	json.NewDecoder(r.Body).Decode(&plan)

	err := validation.Validator.Struct(plan)
	if err != nil {
		validation.ValidatorErrors(w, r, err)
		return
	}

	result := &model.Plan{
		PlanInfo: plan.PlanInfo,
		PlanName: plan.PlanName,
		Status:   plan.Status,
	}

	err = model.DB.Model(&model.Plan{}).Create(&result).Error

	if err != nil {
		log.Fatal(err)
	}

	render.JSON(w, http.StatusCreated, result)
}
