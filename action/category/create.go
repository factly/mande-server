package category

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/factly/x/renderx"
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
// @Success 201 {object} model.Category
// @Failure 400 {array} string
// @Router /categories [post]
func create(w http.ResponseWriter, r *http.Request) {

	category := category{}

	json.NewDecoder(r.Body).Decode(&category)

	validate := validator.New()
	err := validate.Struct(category)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	result := &model.Category{
		Title:    category.Title,
		Slug:     category.Slug,
		ParentID: category.ParentID,
	}

	err = model.DB.Model(&model.Category{}).Create(&result).Error

	if err != nil {
		log.Fatal(err)
	}

	renderx.JSON(w, http.StatusCreated, result)
}
