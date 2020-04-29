package tag

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/models"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
)

// DeleteTag - Delete tag by id
// @Summary Delete a tag
// @Description Delete tag by ID
// @Tags Tag
// @ID delete-tag-by-id
// @Consume  json
// @Param id path string true "Tag ID"
// @Success 200 {object} models.Tag
// @Failure 400 {array} string
// @Router /tags/{id} [delete]
func deleteTag(w http.ResponseWriter, r *http.Request) {

	tagID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(tagID)
	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	tag := &models.Tag{
		ID: uint(id),
	}

	// check record exists or not
	err = models.DB.First(&tag).Error
	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	models.DB.Delete(&tag)

	json.NewEncoder(w).Encode(tag)
}
