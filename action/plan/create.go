package plan

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
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

	req := &model.Plan{}

	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = model.DB.Model(&model.Plan{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(req)
}
