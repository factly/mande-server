package plan

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
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
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	plan := &plan{}
	result := &model.Plan{}
	result.ID = uint(id)

	json.NewDecoder(r.Body).Decode(&plan)

	model.DB.Model(&result).Updates(model.Plan{
		PlanName: plan.PlanName,
		PlanInfo: plan.PlanInfo,
		Status:   plan.Status,
	}).First(&result)

	renderx.JSON(w, http.StatusOK, result)
}
