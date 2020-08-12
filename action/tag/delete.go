package tag

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// delete - Delete tag by id
// @Summary Delete a tag
// @Description Delete tag by ID
// @Tags Tag
// @ID delete-tag-by-id
// @Consume  json
// @Param tag_id path string true "Tag ID"
// @Success 200
// @Failure 400 {array} string
// @Router /tags/{tag_id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {

	tagID := chi.URLParam(r, "tag_id")
	id, err := strconv.Atoi(tagID)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &model.Tag{}
	result.ID = uint(id)

	// check record exists or not
	err = model.DB.First(&result).Error
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	// check if tag is associated with products
	tag := new(model.Tag)
	tag.ID = uint(id)
	totAssociated := model.DB.Model(tag).Association("Products").Count()

	if totAssociated != 0 {
		loggerx.Error(errors.New("tag is associated with product"))
		errorx.Render(w, errorx.Parser(errorx.CannotSaveChanges()))
		return
	}

	model.DB.Delete(&result)

	renderx.JSON(w, http.StatusOK, nil)
}
