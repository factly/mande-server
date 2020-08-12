package medium

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// update - Update medium by id
// @Summary Update a medium by id
// @Description Update medium by ID
// @Tags Medium
// @ID update-medium-by-id
// @Produce json
// @Consume json
// @Param medium_id path string true "Medium ID"
// @Param Medium body medium false "Medium"
// @Success 200 {object} model.Medium
// @Failure 400 {array} string
// @Router /media/{medium_id} [put]
func update(w http.ResponseWriter, r *http.Request) {

	mediumID := chi.URLParam(r, "medium_id")
	id, err := strconv.Atoi(mediumID)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	medium := &medium{}

	json.NewDecoder(r.Body).Decode(&medium)

	result := &model.Medium{}
	result.ID = uint(id)

	model.DB.Model(&result).Update(&model.Medium{
		Name:        medium.Name,
		Slug:        medium.Slug,
		Title:       medium.Title,
		Type:        medium.Type,
		Description: medium.Description,
		Caption:     medium.Caption,
		FileSize:    medium.FileSize,
		URL:         medium.URL,
		Dimensions:  medium.Dimensions,
	}).First(&result)

	renderx.JSON(w, http.StatusOK, result)
}
