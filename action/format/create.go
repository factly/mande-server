package format

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
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
		renderx.JSON(w, http.StatusBadRequest, validationError)
		return
	}

	result := &model.Format{
		Name:        format.Name,
		Description: format.Description,
		IsDefault:   format.IsDefault,
	}

	err := model.DB.Model(&model.Format{}).Create(&result).Error

	if err != nil {
		log.Fatal(err)
	}

	renderx.JSON(w, http.StatusCreated, result)
}
