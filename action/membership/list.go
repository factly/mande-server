package membership

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/middlewarex"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
)

// list response
type paging struct {
	Total int64              `json:"total"`
	Nodes []model.Membership `json:"nodes"`
}

// userList - Get all memberships
// @Summary Show all memberships
// @Description Get all memberships
// @Tags Membership
// @ID get-all-memberships
// @Produce  json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /memberships [get]
func userList(w http.ResponseWriter, r *http.Request) {

	uID, err := middlewarex.GetUser(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := paging{}
	result.Nodes = make([]model.Membership, 0)

	offset, limit := paginationx.Parse(r.URL.Query())

	model.DB.Preload("Plan").Preload("Plan.Catalogs").Preload("Payment").Preload("Payment.Currency").Model(&model.Membership{}).Where(&model.Membership{
		UserID: uint(uID),
	}).Count(&result.Total).Offset(offset).Limit(limit).Find(&result.Nodes)

	renderx.JSON(w, http.StatusOK, result)
}

// adminList - Get all memberships
// @Summary Show all memberships
// @Description Get all memberships
// @Tags Membership
// @ID get-all-memberships
// @Produce  json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param user query string false "User ID"
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /memberships [get]
func adminList(w http.ResponseWriter, r *http.Request) {

	userIDStr := r.URL.Query().Get("user")

	var userID int
	var err error
	if userIDStr != "" {
		userID, err = strconv.Atoi(userIDStr)
		if err != nil {
			loggerx.Error(err)
			errorx.Render(w, errorx.Parser(errorx.InvalidID()))
			return
		}
	}

	result := paging{}
	result.Nodes = make([]model.Membership, 0)

	offset, limit := paginationx.Parse(r.URL.Query())

	tx := model.DB.Model(&model.Membership{}).Preload("Plan").Preload("Plan.Catalogs").Preload("Payment").Preload("Payment.Currency")

	if userID != 0 {
		tx.Where(&model.Membership{
			UserID: uint(userID),
		}).Count(&result.Total).Offset(offset).Limit(limit).Find(&result.Nodes)
	} else {
		tx.Count(&result.Total).Offset(offset).Limit(limit).Find(&result.Nodes)
	}

	renderx.JSON(w, http.StatusOK, result)
}
