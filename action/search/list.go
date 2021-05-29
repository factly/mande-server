package search

import (
	"net/http"
	"strconv"

	"github.com/factly/mande-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/meilisearchx"
	"github.com/factly/x/renderx"
	"github.com/meilisearch/meilisearch-go"
)

// search - Search Entities
// @Summary Global search for all entities
// @Description Global search for all entities
// @Tags Search
// @ID search-entities
// @Produce json
// @Consume json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param q query string true "Query"
// @Param limit query string false "Limit"
// @Success 200
// @Router /search [get]
func list(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query().Get("q")
	if q == "" {
		errorx.Render(w, errorx.Parser(errorx.GetMessage("provide query param q", http.StatusBadRequest)))
		return
	}

	limitQ := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitQ)
	if err != nil && limitQ != "" {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.GetMessage("invalid limit param", http.StatusBadRequest)))
		return
	}

	products := make([]model.Product, 0)
	datasets := make([]model.Dataset, 0)
	catalogs := make([]model.Catalog, 0)

	// search for products
	result, err := meilisearchx.Client.Search("mande").Search(meilisearch.SearchRequest{
		Query:        q,
		Limit:        int64(limit),
		FacetFilters: []string{"kind:product"},
	})

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	productIDs := meilisearchx.GetIDArray(result.Hits)

	// fetch filtered products
	if len(productIDs) > 0 {
		model.DB.Model(&model.Product{}).Where(productIDs).Find(&products)
	}

	// search for datasets
	result, err = meilisearchx.Client.Search("mande").Search(meilisearch.SearchRequest{
		Query:        q,
		Limit:        int64(limit),
		FacetFilters: []string{"kind:dataset"},
	})

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	datasetIDs := meilisearchx.GetIDArray(result.Hits)

	// fetch filtered datasets
	if len(datasetIDs) > 0 {
		model.DB.Model(&model.Dataset{}).Where(datasetIDs).Find(&datasets)
	}

	// search for catalogs
	result, err = meilisearchx.Client.Search("mande").Search(meilisearch.SearchRequest{
		Query:        q,
		Limit:        int64(limit),
		FacetFilters: []string{"kind:catalog"},
	})

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	catalogIDs := meilisearchx.GetIDArray(result.Hits)

	// fetch filtered catalogs
	if len(catalogIDs) > 0 {
		model.DB.Model(&model.Catalog{}).Where(catalogIDs).Find(&catalogs)
	}

	resp := map[string]interface{}{
		"products": products,
		"datasets": datasets,
		"catalogs": catalogs,
	}

	renderx.JSON(w, http.StatusOK, resp)
}
