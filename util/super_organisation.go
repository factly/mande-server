package util

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/mande-server/model"
	"github.com/factly/x/middlewarex"
	"github.com/spf13/viper"
)

// CheckSuperOrganisation checks weather organisation of user is super org or not
func CheckSuperOrganisation(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		oID, err := GetOrganisation(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		uID, err := middlewarex.GetUser(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		superOrgID, err := middlewarex.GetSuperOrganisationID("dataportal")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Check if the organisation id is of super org
		if superOrgID != oID {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Check if the user belong to the organisaiton
		req, err := http.NewRequest("GET", viper.GetString("kavach_url")+"/organisations/my", nil)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User", strconv.Itoa(uID))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		defer resp.Body.Close()

		var orgList []model.OrgWithRole
		err = json.NewDecoder(resp.Body).Decode(&orgList)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		belongs := false
		for _, org := range orgList {
			if org.ID == uint(oID) {
				belongs = true
			}
		}

		if !belongs {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}
