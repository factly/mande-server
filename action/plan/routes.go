package plan

import "github.com/go-chi/chi"

// Plan request body
type plan struct {
	PlanName string `json:"plan_name"`
	PlanInfo string `json:"plan_info"`
	Status   string `json:"status"`
}

// Router - Group of currency router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", getPlans)    // GET /plans - return list of plans
	r.Post("/", CreatePlan) // POST /plans - create a new plan and persist it

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", getPlanByID)   // GET /plans/{id} - read a single plan by :id
		r.Put("/", updatePlan)    // PUT /plans/{id} - update a single plan by :id
		r.Delete("/", deletePlan) // DELETE /plans/{id} - delete a single plan by :id
	})

	return r
}
