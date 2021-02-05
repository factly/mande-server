package user

import (
	"github.com/go-chi/chi"
)

// UserListRouter - Group of tag router
func UserListRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list) // GET /tags - return list of tags

	return r
}
