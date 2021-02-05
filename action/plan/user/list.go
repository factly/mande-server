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
func list(w http.ResponseWriter, r *http.Request) {

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

	result := paging{}

	for _, each := range adminRole.Members {
		id, _ := strconv.Atoi(each)
		result.Nodes = append(result.Nodes, id)
	}

	result.Total = len(adminRole.Members)

	renderx.JSON(w, http.StatusCreated, result)
}
