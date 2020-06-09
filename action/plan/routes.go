package plan

import "github.com/go-chi/chi"

// Plan request body
type plan struct {
	PlanName string `json:"plan_name" validate:"required"`
	PlanInfo string `json:"plan_info"`
	Status   string `json:"status"`
}

// Router - Group of plan router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)    // GET /plans - return list of plans
	r.Post("/", Create) // POST /plans - create a new plan and persist it

	r.Route("/{plan_id}", func(r chi.Router) {
		r.Get("/", details)   // GET /plans/{plan_id} - read a single plan by :plan_id
		r.Put("/", update)    // PUT /plans/{plan_id} - update a single plan by :plan_id
		r.Delete("/", delete) // DELETE /plans/{plan_id} - delete a single plan by :plan_id
	})

	return r
}
