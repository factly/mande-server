package plan

import "github.com/go-chi/chi"

// Plan request body
type plan struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Duration    uint   `json:"duration" validate:"required"`
	Status      string `json:"status"`
	CatalogIDs  []uint `json:"catalog_ids"`
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
