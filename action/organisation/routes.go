package organisation

import "github.com/go-chi/chi"

// UserRouter - Group of plan router
func UserRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list) // GET /organisations - return list of orgs

	return r
}

// AdminRouter - Group of plan router
func AdminRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/my", my) // GET /organisations - return super org

	return r
}
