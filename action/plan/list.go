package plan

import (
	"encoding/json"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
)

// list - Get all plans
// @Summary Show all plans
// @Description Get all plans
// @Tags Plan
// @ID get-all-plans
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} model.Plan
// @Router /plans [get]
func list(w http.ResponseWriter, r *http.Request) {

	var plans []model.Plan

	offset, limit := util.Paging(r.URL.Query())

	model.DB.Offset(offset).Limit(limit).Model(&model.Plan{}).Find(&plans)

	json.NewEncoder(w).Encode(plans)
}
