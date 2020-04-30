package tag

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// update - Update tag by id
// @Summary Update a tag by id
// @Description Update tag by ID
// @Tags Tag
// @ID update-tag-by-id
// @Produce json
// @Consume json
// @Param id path string true "Tag ID"
// @Param Tag body tag false "Tag"
// @Success 200 {object} model.Tag
// @Failure 400 {array} string
// @Router /tags/{id} [put]
func update(w http.ResponseWriter, r *http.Request) {

	tagID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(tagID)
	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &model.Tag{}
	tag := &model.Tag{}
	tag.ID = uint(id)

	json.NewDecoder(r.Body).Decode(&req)

	model.DB.Model(&tag).Update(&model.Tag{Title: req.Title, Slug: req.Slug})
	model.DB.First(&tag)
	json.NewEncoder(w).Encode(tag)
}
