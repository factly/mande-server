package membership

import "github.com/go-chi/chi"

// membership request body
type membership struct {
	PlanID uint `json:"plan_id" validate:"required"`
}

// UserRouter - Group of membership router
func UserRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", userList) // GET /memberships - return list of memberships
	r.Post("/", create)  // POST /memberships - create a new membership and persist it

	r.Route("/{membership_id}", func(r chi.Router) {
		r.Get("/", details)   // GET /memberships/{membership_id} - read a single membership by :membership_id
		r.Delete("/", delete) // DELETE /memberships/{membership_id} - delete a single membership by :membership_id
	})

	return r
}

// AdminRouter - Group of membership router
func AdminRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", adminList) // GET /memberships - return list of memberships

	r.Route("/{membership_id}", func(r chi.Router) {
		r.Get("/", details) // GET /memberships/{membership_id} - read a single membership by :membership_id
	})

	return r
}
