package status

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
)

// GetStatuses - Get all statuses
// @Summary Show all statuses
// @Description Get all statuses
// @Tags Status
// @ID get-all-statuses
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} model.Status
// @Router /products/{id}/status [get]
func GetStatuses(w http.ResponseWriter, r *http.Request) {

	var statuses []model.Status
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

	model.DB.Offset(offset).Limit(limit).Model(&model.Status{}).Find(&statuses)

	json.NewEncoder(w).Encode(statuses)
}
