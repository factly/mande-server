package dataset

import (
	"github.com/factly/mande-server/action/dataset/format"
	"github.com/factly/mande-server/model"
	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm/dialects/postgres"
)

// Dataset request body
type dataset struct {
	Title            string         `json:"title"`
	Description      string         `json:"description"`
	Source           string         `json:"source"`
	Frequency        string         `json:"frequency"`
	TemporalCoverage string         `json:"temporal_coverage"`
	Granularity      string         `json:"granularity"`
	ContactName      string         `json:"contact_name"`
	ContactEmail     string         `json:"contact_email"`
	License          string         `json:"license"`
	DataStandard     string         `json:"data_standard"`
	SampleURL        string         `json:"sample_url"`
	RelatedArticles  postgres.Jsonb `json:"related_articles" swaggertype:"primitive,string"`
	TimeSaved        int            `json:"time_saved"`
	Price            int            `json:"price" validate:"required"`
	CurrencyID       uint           `json:"currency_id"`
	FeaturedMediumID uint           `json:"featured_medium_id"`
	TagIDs           []uint         `json:"tag_ids"`
}

// Dataset detail
type datasetData struct {
	model.Dataset
	Formats []model.DatasetFormat `json:"formats"`
}

var userContext model.ContextKey = "dataset_user"

// PublicRouter - Group of dataset router
func PublicRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", userlist) // GET /datasets - return list of datasets

	r.Route("/{dataset_id}", func(r chi.Router) {
		r.Get("/", userDetails) // GET /datasets/{dataset_id} - read a single dataset by :dataset_id
	})

	return r
}

// AdminRouter - Group of dataset router
func AdminRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", adminlist) // GET /datasets - return list of datasets
	r.Post("/", create)   // POST /datasets - create a new dataset and persist it

	r.Route("/{dataset_id}", func(r chi.Router) {
		r.Get("/", adminDetails) // GET /datasets/{dataset_id} - read a single dataset by :dataset_id
		r.Put("/", update)       // PUT /datasets/{dataset_id} - update a single dataset by :dataset_id
		r.Delete("/", delete)    // DELETE /datasets/{dataset_id} - delete a single dataset by :dataset_id
		r.Mount("/format", format.AdminRouter())
	})

	return r
}
