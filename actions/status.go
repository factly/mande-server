package actions

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/models"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-playground/validator/v10"
)

// status request object
type status struct {
	Name string `json:"name"`
}

// GetStatuses - Get all statuses
// @Summary Show all statuses
// @Description Get all statuses
// @Tags Status
// @ID get-all-statuses
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} models.Status
// @Router /products/{id}/status [get]
func GetStatuses(w http.ResponseWriter, r *http.Request) {

	var statuses []models.Status
	p := r.URL.Query().Get("page")
	pg, _ := strconv.Atoi(p) // pg contains page number
	l := r.URL.Query().Get("limit")
	li, _ := strconv.Atoi(l) // li contains perPage number

	offset := 0 // no. of records to skip
	limit := 5  // limt

	if li > 0 && li <= 10 {
		limit = li
	}

	if pg > 1 {
		offset = (pg - 1) * limit
	}

	models.DB.Offset(offset).Limit(limit).Model(&models.Status{}).Find(&statuses)

	json.NewEncoder(w).Encode(statuses)
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
// @Failure 400 {array} string
// @Router /products/{id}/status [post]
func CreateStatus(w http.ResponseWriter, r *http.Request) {

	req := &models.Status{}
	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = models.DB.Model(&models.Status{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(req)
}
