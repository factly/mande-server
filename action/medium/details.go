package medium

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// details - Get medium by id
// @Summary Show a medium by id
// @Description Get medium by ID
// @Tags Medium
// @ID get-medium-by-id
// @Produce  json
// @Param medium_id path string true "Medium ID"
// @Success 200 {object} model.Medium
// @Failure 400 {array} string
// @Router /media/{medium_id} [get]
func details(w http.ResponseWriter, r *http.Request) {

	mediumID := chi.URLParam(r, "medium_id")
	id, err := strconv.Atoi(mediumID)
	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	result := &model.Medium{}
	result.ID = uint(id)

	err = model.DB.Model(&model.Medium{}).First(&result).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	renderx.JSON(w, http.StatusOK, result)
}
