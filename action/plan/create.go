package plan

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
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

	err := json.NewDecoder(r.Body).Decode(&plan)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	validationError := validationx.Check(plan)
	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	result := &model.Plan{
		PlanInfo: plan.PlanInfo,
		PlanName: plan.PlanName,
		Status:   plan.Status,
	}

	err = model.DB.Model(&model.Plan{}).Create(&result).First(&result).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	renderx.JSON(w, http.StatusCreated, result)
}
