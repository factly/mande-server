package catalog

import (
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
)

// list - Get all catalogs owned by user
// @Summary Show all catalogs owned by user
// @Description Get all catalogs owned by user
// @Tags Product
// @ID get-all-catalogs-owned-by-user
// @Produce  json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /catalogs/my [get]
func my(w http.ResponseWriter, r *http.Request) {

	uID, err := util.GetUser(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &paging{}
	result.Nodes = make([]model.Catalog, 0)

	offset, limit := paginationx.Parse(r.URL.Query())

	memberships := []model.Membership{}

	model.DB.Preload("Plan").Preload("Plan.Catalogs").Preload("Payment").Preload("Payment.Currency").Model(&model.Membership{}).Where(&model.Membership{
		UserID: uint(uID),
	}).Find(&memberships)

	catalogs := []model.Catalog{}

	for _, membership := range memberships {
		catalogs = append(catalogs, membership.Plan.Catalogs...)
	}

	result.Nodes = catalogs[offset : offset+limit]
	result.Total = int64(len(catalogs))

	renderx.JSON(w, http.StatusOK, result)
}
