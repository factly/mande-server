package category

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// update - Update category by id
// @Summary Update a category by id
// @Description Update category by ID
// @Tags Category
// @ID update-category-by-id
// @Produce json
// @Consume json
// @Param category_id path string true "Category ID"
// @Param Category body category false "Category"
// @Success 200 {object} model.Category
// @Failure 400 {array} string
// @Router /categories/{category_id} [put]
func update(w http.ResponseWriter, r *http.Request) {

	categoryID := chi.URLParam(r, "category_id")
	id, err := strconv.Atoi(categoryID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	category := &category{}
	result := &model.Category{}
	result.ID = uint(id)

	json.NewDecoder(r.Body).Decode(&category)

	model.DB.Model(&result).Updates(model.Category{
		Title: category.Title,
		Slug:  category.Slug,
	}).First(&result)

	renderx.JSON(w, http.StatusOK, result)
}
