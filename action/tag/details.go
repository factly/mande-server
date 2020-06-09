package tag

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// details - Get tag by id
// @Summary Show a tag by id
// @Description Get tag by ID
// @Tags Tag
// @ID get-tag-by-id
// @Produce  json
// @Param tag_id path string true "Tag ID"
// @Success 200 {object} model.Tag
// @Failure 400 {array} string
// @Router /tags/{tag_id} [get]
func details(w http.ResponseWriter, r *http.Request) {

	tagID := chi.URLParam(r, "tag_id")
	id, err := strconv.Atoi(tagID)
	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	result := &model.Tag{}
	result.ID = uint(id)

	err = model.DB.Model(&model.Tag{}).First(&result).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	renderx.JSON(w, http.StatusOK, result)
}
