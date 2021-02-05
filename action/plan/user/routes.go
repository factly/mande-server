package user

import (
	"github.com/go-chi/chi"
)

// UserRouter - Group of plan router
func UserRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", list)

	r.Route("/{user_id}", func(r chi.Router) {
		r.Get("/", create)
		r.Delete("/", delete)
	})

	return r
}
