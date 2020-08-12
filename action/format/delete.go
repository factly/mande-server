package format

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
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &model.Format{}

	result.ID = uint(id)

	// check record exists or not
	err = model.DB.First(&result).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	// check if format is associated with datasets
	var totAssociated int
	model.DB.Model(&model.DatasetFormat{}).Where(&model.DatasetFormat{
		FormatID: uint(id),
	}).Count(&totAssociated)

	if totAssociated != 0 {
		loggerx.Error(errors.New("format is associated with dataset"))
		errorx.Render(w, errorx.Parser(errorx.CannotSaveChanges()))
		return
	}

	model.DB.Delete(&result)

	renderx.JSON(w, http.StatusOK, nil)
}
