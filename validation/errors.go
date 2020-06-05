package validation

import (
	"net/http"
	"strings"

	"github.com/factly/x/renderx"
)

// InvalidID - response for invalid ID
func InvalidID(w http.ResponseWriter, r *http.Request) {
	var msg []string
	msg = append(msg, "Invalid id")
	renderx.JSON(w, http.StatusBadRequest, msg)
}

// RecordNotFound - response for record not found
func RecordNotFound(w http.ResponseWriter, r *http.Request) {
	var msg []string
	msg = append(msg, "Record not found")
	renderx.JSON(w, http.StatusNotFound, msg)
}

// ValidErrors - errors from validator
func ValidErrors(w http.ResponseWriter, r *http.Request, msg string) {
	err := strings.Split(msg, "\n")
	renderx.JSON(w, http.StatusBadRequest, err)
}
