package format

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/factly/x/renderx"
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
		validation.InvalidID(w, r)
		return
	}

	format := &format{}
	result := &model.Format{}
	result.ID = uint(id)

	json.NewDecoder(r.Body).Decode(&format)

	model.DB.Model(&result).Updates(model.Format{
		Name:        format.Name,
		Description: format.Description,
		IsDefault:   format.IsDefault,
	}).First(&result)

	renderx.JSON(w, http.StatusOK, result)
}
