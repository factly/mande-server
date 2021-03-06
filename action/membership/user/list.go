package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/data-portal-server/util/keto"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

type paging struct {
	Nodes []int `json:"nodes"`
	Total int   `json:"total"`
}

// list - List Membership users
// @Summary List Membership users
// @Description List Membership users
// @Tags MembershipUser
// @ID list-membership-user
// @Consume json
// @Produce json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param membership_id path string true "Membership ID"
// @Param user_id path string true "User ID"
// @Success 200 {object} paging
// @Failure 400 {array} string
// @Router /memberships/{membership_id}/users [get]
func list(w http.ResponseWriter, r *http.Request) {

	oID, err := util.GetOrganisation(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.Unauthorized()))
		return
	}

	membershipID := chi.URLParam(r, "membership_id")
	memID, err := strconv.Atoi(membershipID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	// Check if membership exist
	membership := model.Membership{}
	membership.ID = uint(memID)

	err = model.DB.First(&membership).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	adminRoleID := fmt.Sprint("roles:org:" + fmt.Sprint(oID) + "app:dataportal:membership:" + fmt.Sprint(memID) + ":users")

	resp, err := keto.GetPolicy("/engines/acp/ory/regex/roles/" + adminRoleID)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.NetworkError()))
		return
	}

	defer resp.Body.Close()

	adminRole := model.Role{}
	err = json.NewDecoder(resp.Body).Decode(&adminRole)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	result := paging{}

	for _, each := range adminRole.Members {
		id, _ := strconv.Atoi(each)
		result.Nodes = append(result.Nodes, id)
	}

	result.Total = len(adminRole.Members)

	renderx.JSON(w, http.StatusOK, result)
}
