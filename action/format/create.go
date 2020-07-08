package format

import (
	"encoding/json"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
)

// create - Create format
// @Summary Create format
// @Description Create format
// @Tags Format
// @ID add-format
// @Consume json
// @Produce  json
// @Param Format body format true "Format object"
// @Success 201 {object} model.Format
// @Failure 400 {array} string
// @Router /formats [post]
func create(w http.ResponseWriter, r *http.Request) {

	format := &format{}
	json.NewDecoder(r.Body).Decode(&format)

	validationError := validationx.Check(format)
	if validationError != nil {
		errorx.Render(w, validationError)
		return
	}

	result := &model.Format{
		Name:        format.Name,
		Description: format.Description,
		IsDefault:   format.IsDefault,
	}

	err := model.DB.Model(&model.Format{}).Create(&result).Error

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	renderx.JSON(w, http.StatusCreated, result)
}
