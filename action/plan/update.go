package plan

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
)

// updatePlan - Update plan by id
// @Summary Update a plan by id
// @Description Update plan by ID
// @Tags Plan
// @ID update-plan-by-id
// @Produce json
// @Consume json
// @Param id path string true "Plan ID"
// @Param Plan body plan false "Plan"
// @Success 200 {object} model.Plan
// @Failure 400 {array} string
// @Router /plans/{id} [put]
func updatePlan(w http.ResponseWriter, r *http.Request) {

	planID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(planID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &model.Plan{}
	plan := &model.Plan{
		ID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&req)

	model.DB.Model(&plan).Updates(model.Plan{
		PlanName: req.PlanName,
		PlanInfo: req.PlanInfo,
		Status:   req.Status,
	})
	model.DB.First(&plan)

	json.NewEncoder(w).Encode(plan)
}
