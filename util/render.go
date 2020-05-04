package util

import (
	"encoding/json"
	"net/http"
)

// Render json
func Render(w http.ResponseWriter, status int, data interface{}) {

	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}
