package category

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-playground/validator/v10"
)

// create - create category
// @Summary Create category
// @Description create category
// @Tags Category
// @ID add-category
// @Consume json
// @Produce  json
// @Param Category body category true "Category object"
// @Success 200 {object} model.Category
// @Failure 400 {array} string
// @Router /categories [post]
func create(w http.ResponseWriter, r *http.Request) {

	req := &model.Category{}

	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = model.DB.Model(&model.Category{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}

	util.Render(w, http.StatusOK, req)
}
