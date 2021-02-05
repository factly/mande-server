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
	"github.com/go-chi/chi"
)

// create - Create organisation
// @Summary Create organisation
// @Description Create organisation
// @Tags Organisation
// @ID add-organisation
// @Consume json
// @Produce json
// @Param X-User header string true "User ID"
// @Param Organisation body organisation true "Organisation Object"
// @Success 201 {object} orgWithRole
// @Failure 400 {array} string
// @Router /organisations [post]
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

	planID := chi.URLParam(r, "plan_id")
	pID, err := strconv.Atoi(planID)

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
	plan := model.Plan{}
	plan.ID = uint(pID)

	err = model.DB.First(&plan).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	adminRoleID := fmt.Sprint("roles:org:" + fmt.Sprint(oID) + ":plan:" + fmt.Sprint(pID) + ":users")

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

	if len(adminRole.Members) >= plan.Users {
		loggerx.Error(errors.New("Cannot add more users"))
		errorx.Render(w, errorx.Parser(errorx.Message{
			Code:    http.StatusUnprocessableEntity,
			Message: "Cannot add more users",
		}))
		return
	}

	/* add user to application */
	reqRole := &model.Role{}
	reqRole.Members = []string{fmt.Sprint(userID)}

	err = keto.UpdateRole("/engines/acp/ory/regex/roles/roles:org:"+fmt.Sprint(oID)+":plan:"+fmt.Sprint(pID)+":users/members", reqRole)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.NetworkError()))
		return
	}

	renderx.JSON(w, http.StatusCreated, nil)
}
