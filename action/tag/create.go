package tag

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
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

	err := validation.Validator.Struct(tag)
	if err != nil {
		validation.ValidatorErrors(w, r, err)
		return
	}
	err = model.DB.Model(&model.Tag{}).Create(&tag).Error

	if err != nil {
		log.Fatal(err)
	}

	render.JSON(w, http.StatusCreated, tag)
}
