package validation

import (
	"encoding/json"
	"net/http"
	"strings"
)

// InvalidID - response for invalid ID
func InvalidID(w http.ResponseWriter, r *http.Request) {
	var msg []string
	msg = append(msg, "Invalid id")
	w.Header().Set("Content-type", "applciation/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(msg)
}

// RecordNotFound - response for record not found
func RecordNotFound(w http.ResponseWriter, r *http.Request) {
	var msg []string
	msg = append(msg, "Record not found")
	w.Header().Set("Content-type", "applciation/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(msg)
}

// ValidErrors - errors from validator
func ValidErrors(w http.ResponseWriter, r *http.Request, msg string) {
	err := strings.Split(msg, "\n")
	w.Header().Set("Content-type", "applciation/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(err)
}
