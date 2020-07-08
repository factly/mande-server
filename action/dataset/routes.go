package dataset

import (
	"github.com/factly/data-portal-server/action/dataset/format"
	"github.com/factly/data-portal-server/model"
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
	RelatedArticles  postgres.Jsonb `json:"related_articles"`
	TimeSaved        int            `json:"time_saved"`
	FeaturedMediumID uint           `json:"featured_medium_id"`
	FormatIDs        []uint         `json:"format_ids"`
	TagIDs           []uint         `json:"tag_ids"`
}

// Dataset detail
type datasetData struct {
	model.Dataset
	Formats []model.DatasetFormat `json:"formats"`
	Tags    []model.Tag           `json:"tags"`
}

// Router - Group of dataset router
func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/", list)    // GET /datasets - return list of datasets
	r.Post("/", create) // POST /datasets - create a new dataset and persist it

	r.Route("/{dataset_id}", func(r chi.Router) {
		r.Get("/", details)   // GET /datasets/{dataset_id} - read a single dataset by :dataset_id
		r.Put("/", update)    // PUT /datasets/{dataset_id} - update a single dataset by :dataset_id
		r.Delete("/", delete) // DELETE /datasets/{dataset_id} - delete a single dataset by :dataset_id
		r.Mount("/format", format.Router())
	})

	return r
}
