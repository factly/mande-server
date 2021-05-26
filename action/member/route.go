package member

import (
	"github.com/go-chi/chi"
)

// Router - Group of member router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)
	return r

}
