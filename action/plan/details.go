package plan

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// details - Get plan by id
// @Summary Show a plan by id
// @Description Get plan by ID
// @Tags Plan
// @ID get-plan-by-id
// @Produce  json
// @Param id path string true "Plan ID"
// @Success 200 {object} model.Plan
// @Failure 400 {array} string
// @Router /plans/{id} [get]
func details(w http.ResponseWriter, r *http.Request) {

	planID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(planID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &model.Plan{}
	req.ID = uint(id)

	err = model.DB.Model(&model.Plan{}).First(&req).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(req)
}
