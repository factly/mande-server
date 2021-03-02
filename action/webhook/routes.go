package webhook

import "github.com/go-chi/chi"

// Router webhooks router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Post("/failed-payment", failedPayment)

	return r
}
