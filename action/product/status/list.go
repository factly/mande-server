package status

import (
	"encoding/json"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
)

// list - Get all statuses
// @Summary Show all statuses
// @Description Get all statuses
// @Tags Status
// @ID get-all-statuses
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} model.Status
// @Router /products/{id}/status [get]
func list(w http.ResponseWriter, r *http.Request) {

	var statuses []model.Status

	offset, limit := util.Paging(r.URL.Query())

	model.DB.Offset(offset).Limit(limit).Model(&model.Status{}).Find(&statuses)

	json.NewEncoder(w).Encode(statuses)
}
