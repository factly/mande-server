package plan

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
)

// deletePlan - Delete plan by id
// @Summary Delete a plan
// @Description Delete plan by ID
// @Tags Plan
// @ID delete-plan-by-id
// @Consume  json
// @Param id path string true "Plan ID"
// @Success 200 {object} model.Plan
// @Failure 400 {array} string
// @Router /plans/{id} [delete]
func deletePlan(w http.ResponseWriter, r *http.Request) {

	planID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(planID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	plan := &model.Plan{
		ID: uint(id),
	}

	// check record exists or not
	err = model.DB.First(&plan).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	model.DB.Delete(&plan)

	json.NewEncoder(w).Encode(plan)
}
