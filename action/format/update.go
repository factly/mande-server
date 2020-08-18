package format

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
	"github.com/go-chi/chi"
)

// update - Update format by id
// @Summary Update a format by id
// @Description Update format by ID
// @Tags Format
// @ID update-format-by-id
// @Produce json
// @Consume json
// @Param format_id path string true "Format ID"
// @Param format body format false "Format"
// @Success 200 {object} model.Format
// @Failure 400 {array} string
// @Router /formats/{format_id} [put]
func update(w http.ResponseWriter, r *http.Request) {

	formatID := chi.URLParam(r, "format_id")
	id, err := strconv.Atoi(formatID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	format := &format{}
	result := &model.Format{}
	result.ID = uint(id)

	json.NewDecoder(r.Body).Decode(&format)

	validationError := validationx.Check(format)
	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	model.DB.Model(&result).Updates(model.Format{
		Name:        format.Name,
		Description: format.Description,
		IsDefault:   format.IsDefault,
	}).First(&result)

	renderx.JSON(w, http.StatusOK, result)
}
