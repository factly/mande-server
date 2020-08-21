package format

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
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
	err := json.NewDecoder(r.Body).Decode(&format)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	validationError := validationx.Check(format)
	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	result := &model.Format{
		Name:        format.Name,
		Description: format.Description,
		IsDefault:   format.IsDefault,
	}

	err = model.DB.Model(&model.Format{}).Create(&result).First(&result).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	renderx.JSON(w, http.StatusCreated, result)
}
