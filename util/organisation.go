package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/factly/mande-server/model"
	"github.com/spf13/viper"
)

type orgWithPermission struct {
	model.Organisation
	Permission permission `json:"permission"`
}

type permission struct {
	Role string `json:"role"`
}

// CheckOwnerFromKavach checks if user is owner of organisation
func CheckOwnerFromKavach(uID, oID int) (bool, error) {
	req, err := http.NewRequest("GET", viper.GetString("kavach_url")+"/organisations/my", nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User", fmt.Sprint(uID))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != http.StatusOK {
		return false, errors.New("error from kavach server")
	}

	var permArr []orgWithPermission

	err = json.NewDecoder(resp.Body).Decode(&permArr)
	if err != nil {
		return false, err
	}

	for _, each := range permArr {
		if each.Permission.Role == "owner" && each.ID == uint(oID) {
			return true, nil
		}
	}

	return false, nil
}
