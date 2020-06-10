package format

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// delete - Delete format by id
// @Summary Delete a format
// @Description Delete format by ID
// @Tags Format
// @ID delete-format-by-id
// @Consume  json
// @Param format_id path string true "Format ID"
// @Success 200
// @Failure 400 {array} string
// @Router /formats/{format_id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {

	formatID := chi.URLParam(r, "format_id")
	id, err := strconv.Atoi(formatID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	result := &model.Format{}

	result.ID = uint(id)

	// check record exists or not
	err = model.DB.First(&result).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	model.DB.Delete(&result)

	renderx.JSON(w, http.StatusOK, nil)
}
