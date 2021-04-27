package catalog

import (
	"net/http"
	"net/url"

	"github.com/factly/mande-server/model"
	"github.com/factly/mande-server/util"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/meilisearchx"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
)

// list response
type paging struct {
	Total int64           `json:"total"`
	Nodes []model.Catalog `json:"nodes"`
}

// list - Get all catalogs
// @Summary Show all catalogs
// @Description Get all catalogs
// @Tags Catalog
// @ID get-all-catalogs
// @Produce  json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /catalogs [get]
func list(w http.ResponseWriter, r *http.Request) {

	result := paging{}
	result.Nodes = make([]model.Catalog, 0)
	offset, limit := paginationx.Parse(r.URL.Query())

	// Filters
	u, _ := url.Parse(r.URL.String())
	queryMap := u.Query()

	searchQuery := r.URL.Query().Get("q")
	sort := r.URL.Query().Get("sort")

	filters := util.GenerateFilters(queryMap["tag"])
	filteredCatalogIDs := make([]uint, 0)
	var err error

	if filters != "" || searchQuery != "" {
		// Search catalog with filter
		var hits []interface{}
		var res map[string]interface{}

		if searchQuery != "" {
			hits, err = meilisearchx.SearchWithQuery("mande", searchQuery, filters, "product")
		} else {
			res, err = meilisearchx.SearchWithoutQuery("mande", filters, "product")
			if _, found := res["hits"]; found {
				hits = res["hits"].([]interface{})
			}
		}
		if err != nil {
			loggerx.Error(err)
			renderx.JSON(w, http.StatusOK, result)
			return
		}

		filteredCatalogIDs = meilisearchx.GetIDArray(hits)
		if len(filteredCatalogIDs) == 0 {
			renderx.JSON(w, http.StatusOK, result)
			return
		}
	}
	if sort != "asc" {
		sort = "desc"
	}

	tx := model.DB.Model(&model.Catalog{}).Preload("FeaturedMedium").Preload("Products").Preload("Products.Currency").Preload("Products.FeaturedMedium").Preload("Products.Tags").Preload("Products.Datasets").Order("created_at " + sort)

	if len(filteredCatalogIDs) > 0 {
		err = tx.Where(filteredCatalogIDs).Count(&result.Total).Offset(offset).Limit(limit).Find(&result.Nodes).Error
	} else {
		err = tx.Count(&result.Total).Offset(offset).Limit(limit).Find(&result.Nodes).Error
	}

	renderx.JSON(w, http.StatusOK, result)
}
