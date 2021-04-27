package tag

import (
	"net/http"

	"github.com/factly/mande-server/model"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/meilisearchx"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
	"log"
)

// list response
type paging struct {
	Total int64       `json:"total"`
	Nodes []model.Tag `json:"nodes"`
}

// @Summary Show all tags
// @Description Get all tags
// @Tags Tag
// @ID get-all-tags
// @Produce  json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /tags [get]
func list(w http.ResponseWriter, r *http.Request) {

	searchQuery := r.URL.Query().Get("q")
	sort := r.URL.Query().Get("sort")

	filteredTagIDs := make([]uint, 0)

	result := paging{}
	result.Nodes = make([]model.Tag, 0)

	var err error

	if searchQuery != "" {

		var hits []interface{}

		hits, err = meilisearchx.SearchWithQuery("mande", searchQuery, "", "tag")

		log.Println("hits", hits)
		log.Println("hits err", err)

		if err != nil {
			loggerx.Error(err)
			renderx.JSON(w, http.StatusOK, result)
			return
		}

		filteredTagIDs = meilisearchx.GetIDArray(hits)
		if len(filteredTagIDs) == 0 {
			renderx.JSON(w, http.StatusOK, result)
			return
		}
	}

	if sort != "asc" {
		sort = "desc"
	}

	offset, limit := paginationx.Parse(r.URL.Query())

	tx := model.DB.Model(&model.Tag{}).Order("created_at " + sort)

	if len(filteredTagIDs) > 0 {
		err = tx.Where(filteredTagIDs).Count(&result.Total).Offset(offset).Limit(limit).Find(&result.Nodes).Error
	} else {
		err = tx.Count(&result.Total).Offset(offset).Limit(limit).Find(&result.Nodes).Error
	}

	renderx.JSON(w, http.StatusOK, result)
}
