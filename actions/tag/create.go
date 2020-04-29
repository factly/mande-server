package tag

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-api/models"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-playground/validator/v10"
)

// CreateTag - Create tag
// @Summary Create tag
// @Description Create tag
// @Tags Tag
// @ID add-tag
// @Consume json
// @Produce  json
// @Param Tag body tag true "Tag object"
// @Success 200 {object} models.Tag
// @Failure 400 {array} string
// @Router /tags [post]
func createTag(w http.ResponseWriter, r *http.Request) {

	req := &models.Tag{}

	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}
	err = models.DB.Model(&models.Tag{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(req)
}
