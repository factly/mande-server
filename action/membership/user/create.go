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
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
	"github.com/go-chi/chi"
)

// create - Create Membership user
// @Summary Create Membership user
// @Description Create Membership user
// @Tags MembershipUser
// @ID create-membership-user
// @Consume json
// @Produce json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param membership_id path string true "Membership ID"
// @Param Request body userRequest true "User Request Object"
// @Success 201
// @Failure 400 {array} string
// @Router /memberships/{membership_id}/users [post]
func create(w http.ResponseWriter, r *http.Request) {
	uID, err := util.GetUser(r.Context())
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

	user := &userRequest{}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	validationError := validationx.Check(user)
	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
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

	err = model.DB.Preload("Plan").First(&membership).Error

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

	if len(adminRole.Members) >= membership.Plan.Users {
		loggerx.Error(errors.New("Cannot add more users"))
		errorx.Render(w, errorx.Parser(errorx.Message{
			Code:    http.StatusUnprocessableEntity,
			Message: "Cannot add more users",
		}))
		return
	}

	/* add user to application */
	reqRole := &model.Role{}
	reqRole.Members = []string{fmt.Sprint(user.UserID)}

	err = keto.UpdateRole("/engines/acp/ory/regex/roles/roles:org:"+fmt.Sprint(oID)+"app:dataportal:membership:"+fmt.Sprint(memID)+":users/members", reqRole)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.NetworkError()))
		return
	}

	renderx.JSON(w, http.StatusCreated, nil)
}
