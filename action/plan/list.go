package plan

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
)

// getPlans - Get all plans
// @Summary Show all plans
// @Description Get all plans
// @Tags Plan
// @ID get-all-plans
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} model.Plan
// @Router /plans [get]
func getPlans(w http.ResponseWriter, r *http.Request) {

	var plans []model.Plan
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

	model.DB.Offset(offset).Limit(limit).Model(&model.Plan{}).Find(&plans)

	json.NewEncoder(w).Encode(plans)
}
