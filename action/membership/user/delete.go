package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/data-portal-server/util/keto"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/middlewarex"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// delete - Delete Membership user
// @Summary Delete Membership user
// @Description Delete Membership user
// @Tags MembershipUser
// @ID delete-membership-user
// @Consume json
// @Produce json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param membership_id path string true "Membership ID"
// @Param user_id path string true "User ID"
// @Success 200
// @Failure 400 {array} string
// @Router /memberships/{membership_id}/users/{user_id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {
	uID, err := middlewarex.GetUser(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.Unauthorized()))
		return
	}

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

	planUserID := chi.URLParam(r, "user_id")
	userID, err := strconv.Atoi(planUserID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	// Check if logged in user is owner
	isAdmin, err := util.CheckOwnerFromKavach(uID, oID)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	if !isAdmin {
		loggerx.Error(errors.New("user is not admin"))
		errorx.Render(w, errorx.Parser(errorx.Unauthorized()))
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

	if len(adminRole.Members) < 2 {
		loggerx.Error(errors.New("Cannot add delete last user"))
		errorx.Render(w, errorx.Parser(errorx.Message{
			Code:    http.StatusUnprocessableEntity,
			Message: "Cannot add more user",
		}))
		return
	}

	err = keto.DeletePolicy("roles:org:" + fmt.Sprint(oID) + "app:dataportal:membership:" + fmt.Sprint(memID) + ":users/members/" + fmt.Sprint(userID))

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.NetworkError()))
		return
	}

	renderx.JSON(w, http.StatusOK, nil)
}
