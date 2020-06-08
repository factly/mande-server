package membership

import "github.com/go-chi/chi"

// membership request body
type membership struct {
	Status    string `json:"status"`
	UserID    uint   `json:"user_id" validate:"required"`
	PaymentID uint   `json:"payment_id" validate:"required"`
	PlanID    uint   `json:"plan_id" validate:"required"`
}

// Router - Group of membership router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)    // GET /memberships - return list of memberships
	r.Post("/", create) // POST /memberships - create a new membership and persist it

	r.Route("/{membership_id}", func(r chi.Router) {
		r.Get("/", details)   // GET /memberships/{membership_id} - read a single membership by :membership_id
		r.Put("/", update)    // PUT /memberships/{membership_id} - update a single membership by :membership_id
		r.Delete("/", delete) // DELETE /memberships/{membership_id} - delete a single membership by :membership_id
	})

	return r
}
