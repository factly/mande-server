package membership

import "github.com/go-chi/chi"

// membership request body
type membership struct {
	Status    string `json:"status"`
	UserID    uint   `json:"user_id"`
	PaymentID uint   `json:"payment_id"`
	PlanID    uint   `json:"plan_id"`
}

// Router - Group of membership router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)    // GET /memberships - return list of memberships
	r.Post("/", create) // POST /memberships - create a new membership and persist it

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", details)   // GET /memberships/{id} - read a single membership by :id
		r.Put("/", update)    // PUT /memberships/{id} - update a single membership by :id
		r.Delete("/", delete) // DELETE /memberships/{id} - delete a single membership by :id
	})

	return r
}
