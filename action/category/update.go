package category

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
)

// updateCategory - Update category by id
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
func updateCategory(w http.ResponseWriter, r *http.Request) {

	categoryID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(categoryID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &model.Category{}
	category := &model.Category{
		ID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&req)

	model.DB.Model(&category).Updates(model.Category{
		Title: req.Title,
		Slug:  req.Slug,
	})
	model.DB.First(&category)

	json.NewEncoder(w).Encode(category)
}
