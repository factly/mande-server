package tag

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
)

// deleteTag - Delete tag by id
// @Summary Delete a tag
// @Description Delete tag by ID
// @Tags Tag
// @ID delete-tag-by-id
// @Consume  json
// @Param id path string true "Tag ID"
// @Success 200 {object} model.Tag
// @Failure 400 {array} string
// @Router /tags/{id} [delete]
func deleteTag(w http.ResponseWriter, r *http.Request) {

	tagID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(tagID)
	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	tag := &model.Tag{
		ID: uint(id),
	}

	// check record exists or not
	err = model.DB.First(&tag).Error
	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	model.DB.Delete(&tag)

	json.NewEncoder(w).Encode(tag)
}
