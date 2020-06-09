package plan

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// delete - Delete plan by id
// @Summary Delete a plan
// @Description Delete plan by ID
// @Tags Plan
// @ID delete-plan-by-id
// @Consume  json
// @Param plan_id path string true "Plan ID"
// @Success 200
// @Failure 400 {array} string
// @Router /plans/{plan_id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {

	planID := chi.URLParam(r, "plan_id")
	id, err := strconv.Atoi(planID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	result := &model.Plan{}

	result.ID = uint(id)

	// check record exists or not
	err = model.DB.First(&result).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	model.DB.Delete(&result)

	renderx.JSON(w, http.StatusOK, nil)
}
