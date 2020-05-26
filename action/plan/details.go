package plan

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// details - Get plan by id
// @Summary Show a plan by id
// @Description Get plan by ID
// @Tags Plan
// @ID get-plan-by-id
// @Produce  json
// @Param plan_id path string true "Plan ID"
// @Success 200 {object} model.Plan
// @Failure 400 {array} string
// @Router /plans/{plan_id} [get]
func details(w http.ResponseWriter, r *http.Request) {

	planID := chi.URLParam(r, "plan_id")
	id, err := strconv.Atoi(planID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	result := &model.Plan{}
	result.ID = uint(id)

	err = model.DB.Model(&model.Plan{}).First(&result).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	render.JSON(w, http.StatusOK, result)
}
