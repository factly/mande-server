package status

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-playground/validator/v10"
)

// create - Create status
// @Summary Create status
// @Description Create status
// @Tags Status
// @ID add-status
// @Consume json
// @Produce  json
// @Param Status body status true "Status object"
// @Success 201 {object} model.Status
// @Failure 400 {array} string
// @Router /products/{id}/status [post]
func create(w http.ResponseWriter, r *http.Request) {

	productStatus := &model.Status{}
	json.NewDecoder(r.Body).Decode(&productStatus)

	validate := validator.New()
	err := validate.Struct(productStatus)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = model.DB.Model(&model.Status{}).Create(&productStatus).Error

	if err != nil {
		log.Fatal(err)
	}

	render.JSON(w, http.StatusCreated, productStatus)
}
