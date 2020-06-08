package tag

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
)

// create - Create tag
// @Summary Create tag
// @Description Create tag
// @Tags Tag
// @ID add-tag
// @Consume json
// @Produce  json
// @Param Tag body tag true "Tag object"
// @Success 201 {object} model.Tag
// @Failure 400 {array} string
// @Router /tags [post]
func create(w http.ResponseWriter, r *http.Request) {

	tag := &model.Tag{}

	json.NewDecoder(r.Body).Decode(&tag)

	validationError := validationx.Check(tag)
	if validationError != nil {
		renderx.JSON(w, http.StatusBadRequest, validationError)
		return
	}

	err := model.DB.Model(&model.Tag{}).Create(&tag).Error

	if err != nil {
		log.Fatal(err)
	}

	renderx.JSON(w, http.StatusCreated, tag)
}
