package organisation

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/factly/mande-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/middlewarex"
	"github.com/factly/x/renderx"
	"github.com/spf13/viper"
)

type organisationApplication struct {
	model.Base
	Title        string        `gorm:"column:title" json:"title"`
	Slug         string        `gorm:"column:slug;unique_index" json:"slug"`
	Applications []application `json:"applications"`
}

type application struct {
	model.Base
	Name        string        `gorm:"column:name" json:"name"`
	Description string        `gorm:"column:description" json:"description"`
	URL         string        `gorm:"column:url" json:"url"`
	MediumID    *uint         `gorm:"column:medium_id;default:NULL" json:"medium_id"`
	Medium      *model.Medium `gorm:"foreignKey:medium_id" json:"medium"`
}

// my - Get super org
// @Summary  Get super org
// @Description  Get super org
// @Tags Organisation
// @ID get-super-organisation
// @Produce  json
// @Param X-User header string true "User ID"
// @Success 200 {object} organisationApplication
// @Router /admin/organisations/my [get]
func my(w http.ResponseWriter, r *http.Request) {
	uID, err := middlewarex.GetUser(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.Unauthorized()))
		return
	}

	// Fetched all organisations of the user
	req, err := http.NewRequest("GET", viper.GetString("kavach_url")+"/organisations/my", nil)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}
	req.Header.Set("X-User", strconv.Itoa(uID))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.NetworkError()))
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	allOrg := []organisationApplication{}
	err = json.Unmarshal(body, &allOrg)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	superOrgID, err := middlewarex.GetSuperOrganisationID("mande")
	if err != nil {
		renderx.JSON(w, http.StatusUnauthorized, nil)
		return
	}

	for _, org := range allOrg {
		if org.ID == uint(superOrgID) {
			renderx.JSON(w, http.StatusOK, org)
		}
	}

	renderx.JSON(w, http.StatusUnauthorized, nil)
}

// list - Get all org
// @Summary  Get all org
// @Description  Get all org
// @Tags Organisation
// @ID get-all-organisation
// @Produce  json
// @Param X-User header string true "User ID"
// @Success 200 {array} organisationApplication
// @Router /admin/organisations/my [get]
func list(w http.ResponseWriter, r *http.Request) {
	uID, err := middlewarex.GetUser(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.Unauthorized()))
		return
	}

	// Fetched all organisations of the user
	req, err := http.NewRequest("GET", viper.GetString("kavach_url")+"/organisations/my", nil)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}
	req.Header.Set("X-User", strconv.Itoa(uID))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.NetworkError()))
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	allOrg := []organisationApplication{}
	err = json.Unmarshal(body, &allOrg)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	renderx.JSON(w, http.StatusUnauthorized, allOrg)
}
