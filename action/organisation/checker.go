package organisation

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/spf13/viper"
)

type authenticationSession struct {
	Subject      string                 `json:"subject"`
	Extra        map[string]interface{} `json:"extra"`
	Header       http.Header            `json:"header"`
	MatchContext matchContext           `json:"match_context"`
}

type matchContext struct {
	RegexpCaptureGroups []string `json:"regexp_capture_groups"`
	URL                 *url.URL `json:"url"`
}

func checker(w http.ResponseWriter, r *http.Request) {
	oID, err := util.GetOrganisation(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	uID, err := util.GetUser(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	payload := &authenticationSession{}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	superOrgID, err := util.GetSuperOrganisationID()
	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	// Check if the organisation id is of super org
	if superOrgID != oID {
		errorx.Render(w, errorx.Parser(errorx.GetMessage("Not Super Organisation", http.StatusUnauthorized)))
		return
	}

	// Check if the user belong to the organisaiton
	req, err := http.NewRequest("GET", viper.GetString("kavach_url")+"/organisations/my", nil)
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

	var orgList []model.OrgWithRole
	err = json.NewDecoder(resp.Body).Decode(&orgList)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	belongs := false
	for _, org := range orgList {
		if org.ID == uint(oID) {
			belongs = true
		}
	}

	if belongs {
		errorx.Render(w, errorx.Parser(errorx.GetMessage("User is not from Super Organisation", http.StatusUnauthorized)))
		return
	}

	renderx.JSON(w, http.StatusOK, payload)
}
