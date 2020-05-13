package category

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// update - Update category by id
// @Summary Update a category by id
// @Description Update category by ID
// @Tags Category
// @ID update-category-by-id
// @Produce json
// @Consume json
// @Param id path string true "Category ID"
// @Param Category body category false "Category"
// @Success 200 {object} model.Category
// @Failure 400 {array} string
// @Router /categories/{id} [put]
func update(w http.ResponseWriter, r *http.Request) {

	categoryID := chi.URLParam(r, "id")
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
	})
	model.DB.First(&result)

	render.JSON(w, http.StatusOK, result)
}
