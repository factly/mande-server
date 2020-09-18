package util

import (
	"errors"
	"net/http"
	"strconv"
)

// GetUser return user ID
func GetUser(r *http.Request) (int, error) {
	user := r.Header.Get("X-User")
	if user == "" {
		return 0, errors.New("no userID")
	}

	uID, err := strconv.Atoi(user)
	if err != nil {
		return 0, err
	}
	return uID, nil
}
