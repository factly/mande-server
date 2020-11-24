package util

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/spf13/viper"
)

// GetSuperOrganisationID get superorganisation id from keto policy
func GetSuperOrganisationID() (int, error) {

	req, err := http.NewRequest("GET", viper.GetString("keto_url")+"/engines/acp/ory/regex/policies/app:dataportal:superorg", nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode == http.StatusOK {
		var policy model.KetoPolicy
		err = json.NewDecoder(resp.Body).Decode(&policy)
		if err != nil {
			return 0, err
		}

		if len(policy.Subjects) != 0 {
			orgID, _ := strconv.Atoi(policy.Subjects[0])
			return orgID, nil
		}
	}
	return 0, errors.New("cannot get super organisation id")
}

// CheckSuperOrganisation checks weather organisation of user is super org or not
func CheckSuperOrganisation(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		oID, err := GetOrganisation(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		uID, err := GetUser(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		superOrgID, err := GetSuperOrganisationID()
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
