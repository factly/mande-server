package tag

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
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
// @Param tag_id path string true "Tag ID"
// @Param Tag body tag false "Tag"
// @Success 200 {object} model.Tag
// @Failure 400 {array} string
// @Router /tags/{tag_id} [put]
func update(w http.ResponseWriter, r *http.Request) {

	tagID := chi.URLParam(r, "tag_id")
	id, err := strconv.Atoi(tagID)
	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	tag := &tag{}

	json.NewDecoder(r.Body).Decode(&tag)

	result := &model.Tag{}
	result.ID = uint(id)

	model.DB.Model(&result).Update(&model.Tag{
		Title: tag.Title,
		Slug:  tag.Slug,
	}).First(&result)

	render.JSON(w, http.StatusOK, result)
}
