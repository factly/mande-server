package membership

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
)

// getMemberships - Get all memberships
// @Summary Show all memberships
// @Description Get all memberships
// @Tags Membership
// @ID get-all-memberships
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} model.Membership
// @Router /memberships [get]
func getMemberships(w http.ResponseWriter, r *http.Request) {

	var memberships []model.Membership
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

	model.DB.Offset(offset).Limit(limit).Preload("User").Preload("Plan").Preload("Payment").Preload("Payment.Currency").Model(&model.Membership{}).Find(&memberships)

	json.NewEncoder(w).Encode(memberships)
}
