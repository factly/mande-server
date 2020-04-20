package actions

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-api/models"
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
// @Success 200 {object} models.Status
// @Router /products/{id}/status [post]
func CreateStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	req := &models.Status{}
	json.NewDecoder(r.Body).Decode(&req)

	if validErrs := req.Validate(); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		w.Header().Set("Content-type", "applciation/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	err := models.DB.Model(&models.Status{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(req)
}
