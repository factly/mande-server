package format

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// details - Get format by id
// @Summary Show a format by id
// @Description Get format by ID
// @Tags Format
// @ID get-format-by-id
// @Produce  json
// @Param format_id path string true "format ID"
// @Success 200 {object} model.Format
// @Failure 400 {array} string
// @Router /formats/{format_id} [get]
func details(w http.ResponseWriter, r *http.Request) {

	formatID := chi.URLParam(r, "format_id")
	id, err := strconv.Atoi(formatID)

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &model.Format{}
	result.ID = uint(id)

	err = model.DB.Model(&model.Format{}).First(&result).Error

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	renderx.JSON(w, http.StatusOK, result)
}
