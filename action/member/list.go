package member

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/factly/mande-server/model"
	"github.com/factly/mande-server/util"
	"github.com/factly/x/loggerx"
	"github.com/spf13/viper"

	"github.com/factly/x/errorx"
	"github.com/factly/x/middlewarex"
	"github.com/factly/x/renderx"
)

// list response
type paging struct {
	Total int            `json:"total"`
	Nodes []model.Member `json:"nodes"`
}

// list - Get all members
// @Summary Show all members
// @Description Get all members
// @Tags Members
// @ID get-all-members
// @Produce  json
// @Param X-User header string true "User ID"
// @Param X-Space header string true "Space ID"
// @Param limit query string false "limit per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /core/members [get]
func list(w http.ResponseWriter, r *http.Request) {

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

	result := paging{}
	result.Nodes = make([]model.Member, 0)

	url := fmt.Sprint(viper.GetString("kavach_url"), "/organisations/", oID, "/users")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User", strconv.Itoa(uID))
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.NetworkError()))
		return
	}

	defer resp.Body.Close()

	users := make([]model.Member, 0)
	err = json.NewDecoder(resp.Body).Decode(&users)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	total := len(users)

	result.Nodes = users
	result.Total = total

	renderx.JSON(w, http.StatusOK, result)
}
