package user

import (
	"github.com/go-chi/chi"
)

type userRequest struct {
	UserID uint `json:"user_id" validate:"required"`
}

// UserRouter - Group of plan router
func UserRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", list)
	r.Post("/", create)

	r.Route("/{user_id}", func(r chi.Router) {
		r.Delete("/", delete)
	})

	return r
}
