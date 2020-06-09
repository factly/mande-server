package user

import "github.com/go-chi/chi"

// user request body
type user struct {
	Email     string `json:"email" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"`
}

// Router - Group of user router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)    // GET /users - return list of users
	r.Post("/", create) // POST /users - create a new user and persist it

	r.Route("/{user_id}", func(r chi.Router) {
		r.Get("/", details)   // GET /users/{user_id} - read a single user by :user_id
		r.Put("/", update)    // PUT /users/{user_id} - update a single user by :user_id
		r.Delete("/", delete) // DELETE /users/{user_id} - delete a single user by :user_id
	})

	return r
}
