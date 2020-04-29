package status

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-playground/validator/v10"
)

// status request object
type status struct {
	Name string `json:"name"`
}

// CreateStatus - Create status
// @Summary Create status
// @Description Create status
// @Tags Status
// @ID add-status
// @Consume json
// @Produce  json
// @Param Status body status true "Status object"
// @Success 200 {object} model.Status
// @Failure 400 {array} string
// @Router /products/{id}/status [post]
func CreateStatus(w http.ResponseWriter, r *http.Request) {

	req := &model.Status{}
	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = model.DB.Model(&model.Status{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(req)
}
