package medium

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
)

// create - Create medium
// @Summary Create medium
// @Description Create medium
// @Tags Medium
// @ID add-medium
// @Consume json
// @Produce  json
// @Param Medium body medium true "Medium object"
// @Success 201 {object} model.Medium
// @Failure 400 {array} string
// @Router /media [post]
func create(w http.ResponseWriter, r *http.Request) {

	medium := &medium{}

	json.NewDecoder(r.Body).Decode(&medium)

	validationError := validationx.Check(medium)
	if validationError != nil {
		renderx.JSON(w, http.StatusBadRequest, validationError)
		return
	}

	result := &model.Medium{
		Name:        medium.Name,
		Slug:        medium.Slug,
		Title:       medium.Title,
		Type:        medium.Type,
		Description: medium.Description,
		Caption:     medium.Caption,
		FileSize:    medium.FileSize,
		URL:         medium.URL,
		Dimensions:  medium.Dimensions,
	}

	err := model.DB.Model(&model.Medium{}).Create(&result).Error

	if err != nil {
		log.Fatal(err)
	}

	renderx.JSON(w, http.StatusCreated, result)
}
