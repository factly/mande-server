package medium

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/factly/mande-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/meilisearchx"
	"github.com/factly/x/middlewarex"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
)

// create - Create medium
// @Summary Create medium
// @Description Create medium
// @Tags Medium
// @ID add-medium
// @Consume json
// @Produce  json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param Medium body medium true "Medium object"
// @Success 201 {object} model.Medium
// @Failure 400 {array} string
// @Router /media [post]
func create(w http.ResponseWriter, r *http.Request) {
	uID, err := middlewarex.GetUser(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	medium := &medium{}

	err = json.NewDecoder(r.Body).Decode(&medium)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	validationError := validationx.Check(medium)
	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	result := &model.Medium{
		Name:        medium.Name,
		Slug:        medium.Slug,
		Title:       medium.Title,
		Type:        medium.Type,
		Description: medium.Description,
		Caption:     medium.Caption,
		AltText:     medium.AltText,
		FileSize:    medium.FileSize,
		URL:         medium.URL,
		Dimensions:  medium.Dimensions,
	}

	tx := model.DB.WithContext(context.WithValue(r.Context(), userContext, uID)).Begin()
	err = tx.Model(&model.Medium{}).Create(&result).Error

	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}
	// Insert into meili index
	meiliObj := map[string]interface{}{
		"id":          result.ID,
		"kind":        "medium",
		"name":        result.Name,
		"slug":        result.Slug,
		"type":        result.Type,
		"title":       result.Title,
		"description": result.Description,
		"caption":     result.Caption,
		"alt_text":    result.AltText,
		"file_size":   result.FileSize,
		"dimensions":  result.Dimensions,
	}

	err = meilisearchx.AddDocument("data-portal", meiliObj)
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	tx.Commit()
	renderx.JSON(w, http.StatusCreated, result)
}
