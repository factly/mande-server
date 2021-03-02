package tag

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/meilisearchx"
	"github.com/factly/x/middlewarex"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
	"github.com/go-chi/chi"
)

// update - Update tag by id
// @Summary Update a tag by id
// @Description Update tag by ID
// @Tags Tag
// @ID update-tag-by-id
// @Produce json
// @Consume json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param tag_id path string true "Tag ID"
// @Param Tag body tag false "Tag"
// @Success 200 {object} model.Tag
// @Failure 400 {array} string
// @Router /tags/{tag_id} [put]
func update(w http.ResponseWriter, r *http.Request) {
	uID, err := middlewarex.GetUser(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.Unauthorized()))
		return
	}

	tagID := chi.URLParam(r, "tag_id")
	id, err := strconv.Atoi(tagID)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	tag := &tag{}

	err = json.NewDecoder(r.Body).Decode(&tag)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	validationError := validationx.Check(tag)
	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	result := &model.Tag{}
	result.ID = uint(id)

	err = model.DB.First(&result).Error
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	tx := model.DB.Begin()
	tx.Model(&result).Updates(&model.Tag{
		Base:  model.Base{UpdatedByID: uint(uID)},
		Title: tag.Title,
		Slug:  tag.Slug,
	}).First(&result)

	// Update into meili index
	meiliObj := map[string]interface{}{
		"id":    result.ID,
		"kind":  "tag",
		"title": result.Title,
		"slug":  result.Slug,
	}

	err = meilisearchx.UpdateDocument("data-portal", meiliObj)
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	tx.Commit()
	renderx.JSON(w, http.StatusOK, result)
}
