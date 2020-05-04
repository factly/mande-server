package category

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// details - Get category by id
// @Summary Show a category by id
// @Description Get category by ID
// @Tags Category
// @ID get-category-by-id
// @Produce  json
// @Param id path string true "Category ID"
// @Success 200 {object} model.Category
// @Failure 400 {array} string
// @Router /categories/{id} [get]
func details(w http.ResponseWriter, r *http.Request) {

	categoryID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(categoryID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	category := &model.Category{}
	category.ID = uint(id)

	err = model.DB.Model(&model.Category{}).First(&category).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	util.Render(w, http.StatusOK, category)
}
