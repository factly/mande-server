package validation

import (
	"net/http"
	"strings"

	"github.com/factly/data-portal-server/util/render"
	"github.com/go-playground/validator/v10"
)

// InvalidID - response for invalid ID
func InvalidID(w http.ResponseWriter, r *http.Request) {
	var msg []string
	msg = append(msg, "Invalid id")
	render.JSON(w, http.StatusBadRequest, msg)
}

// RecordNotFound - response for record not found
func RecordNotFound(w http.ResponseWriter, r *http.Request) {
	var msg []string
	msg = append(msg, "Record not found")
	render.JSON(w, http.StatusNotFound, msg)
}

// ValidErrors - errors from validator
func ValidErrors(w http.ResponseWriter, r *http.Request, msg string) {
	err := strings.Split(msg, "\n")
	render.JSON(w, http.StatusBadRequest, err)
}

// ValidatorErrors - errors from validator
func ValidatorErrors(w http.ResponseWriter, r *http.Request, err error) {
	msg := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		msg[e.Field()] = e.Translate(Trans)
	}
	render.JSON(w, http.StatusBadRequest, msg)
}
