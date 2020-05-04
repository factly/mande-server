package plan

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-playground/validator/v10"
)

// Create - create plan
// @Summary Create plan
// @Description create plan
// @Tags Plan
// @ID add-plan
// @Consume json
// @Produce  json
// @Param Plan body plan true "Plan object"
// @Success 200 {object} model.Plan
// @Router /plans [post]
func Create(w http.ResponseWriter, r *http.Request) {

	plan := &model.Plan{}

	json.NewDecoder(r.Body).Decode(&plan)

	validate := validator.New()
	err := validate.Struct(plan)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = model.DB.Model(&model.Plan{}).Create(&plan).Error

	if err != nil {
		log.Fatal(err)
	}

	util.Render(w, http.StatusOK, plan)
}
