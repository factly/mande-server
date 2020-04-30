package user

import "github.com/go-chi/chi"

// user request body
type user struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Router - Group of currency router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)    // GET /users - return list of users
	r.Post("/", create) // POST /users - create a new user and persist it

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", details)   // GET /users/{id} - read a single user by :id
		r.Put("/", update)    // PUT /users/{id} - update a single user by :id
		r.Delete("/", delete) // DELETE /users/{id} - delete a single user by :id
	})

	return r
}
