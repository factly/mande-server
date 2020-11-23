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
