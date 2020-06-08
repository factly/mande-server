package category

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
)

// create - create category
// @Summary Create category
// @Description create category
// @Tags Category
// @ID add-category
// @Consume json
// @Produce  json
// @Param Category body category true "Category object"
// @Success 201 {object} model.Category
// @Failure 400 {array} string
// @Router /categories [post]
func create(w http.ResponseWriter, r *http.Request) {

	category := category{}

	json.NewDecoder(r.Body).Decode(&category)

	validationError := validationx.Check(category)
	if validationError != nil {
		validation.ValidatorErrors(w, r, validationError)
		return
	}

	result := &model.Category{
		Title:    category.Title,
		Slug:     category.Slug,
		ParentID: category.ParentID,
	}

	err := model.DB.Model(&model.Category{}).Create(&result).Error

	if err != nil {
		log.Fatal(err)
	}

	renderx.JSON(w, http.StatusCreated, result)
}
