package plan

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// update - Update plan by id
// @Summary Update a plan by id
// @Description Update plan by ID
// @Tags Plan
// @ID update-plan-by-id
// @Produce json
// @Consume json
// @Param plan_id path string true "Plan ID"
// @Param Plan body plan false "Plan"
// @Success 200 {object} model.Plan
// @Failure 400 {array} string
// @Router /plans/{plan_id} [put]
func update(w http.ResponseWriter, r *http.Request) {

	planID := chi.URLParam(r, "plan_id")
	id, err := strconv.Atoi(planID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &model.Plan{}
	plan := &model.Plan{}
	plan.ID = uint(id)

	json.NewDecoder(r.Body).Decode(&req)

	model.DB.Model(&plan).Updates(model.Plan{
		PlanName: req.PlanName,
		PlanInfo: req.PlanInfo,
		Status:   req.Status,
	})
	model.DB.First(&plan)

	render.JSON(w, http.StatusOK, plan)
}
