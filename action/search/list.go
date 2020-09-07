package search

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/factly/data-portal-server/util/meili"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
	"github.com/meilisearch/meilisearch-go"
)

// search - Search Entities
// @Summary Global search for all entities
// @Description Global search for all entities
// @Tags Search
// @ID search-entities
// @Produce json
// @Consume json
// @Param Search body searchQuery false "Search"
// @Success 200
// @Router /search [post]
func list(w http.ResponseWriter, r *http.Request) {

	searchQuery := &searchQuery{}
	err := json.NewDecoder(r.Body).Decode(&searchQuery)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	validationError := validationx.Check(searchQuery)
	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	result, err := meili.Client.Search("data-portal").Search(meilisearch.SearchRequest{
		Query:        searchQuery.Query,
		Limit:        searchQuery.Limit,
		Filters:      searchQuery.Filters,
		FacetFilters: searchQuery.FacetFilters,
	})

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	renderx.JSON(w, http.StatusOK, result.Hits)
}
