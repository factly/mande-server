package dataset

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
	Total int64         `json:"total"`
	Nodes []datasetData `json:"nodes"`
}

// userlist - Get all datsets
// @Summary Show all datsets
// @Description Get all datsets
// @Tags Dataset
// @ID get-all-datsets
// @Produce  json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /datasets [get]
func userlist(w http.ResponseWriter, r *http.Request) {
	nodes := make([]datasetData, 0)
	result := paging{}
	result.Nodes = make([]datasetData, 0)
	datasets := make([]model.Dataset, 0)

	offset, limit := paginationx.Parse(r.URL.Query())

	model.DB.Preload("FeaturedMedium").Preload("Currency").Preload("Tags").Model(&model.Dataset{}).Count(&result.Total).Offset(offset).Limit(limit).Find(&datasets)

	for _, dataset := range datasets {
		var formats []model.DatasetFormat

		data := &datasetData{}
		data.Formats = make([]model.DatasetFormat, 0)
		model.DB.Model(&model.DatasetFormat{}).Select("id, created_at, updated_at, deleted_at, format_id, dataset_id").Where(&model.DatasetFormat{
			DatasetID: uint(dataset.ID),
		}).Preload("Format").Find(&formats)

		data.Formats = append(data.Formats, formats...)

		data.Dataset = dataset

		nodes = append(nodes, *data)
	}

	result.Nodes = nodes

	renderx.JSON(w, http.StatusOK, result)
}

// adminlist - Get all datsets
// @Summary Show all datsets
// @Description Get all datsets
// @Tags Dataset
// @ID get-all-datsets
// @Produce  json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /datasets [get]
func adminlist(w http.ResponseWriter, r *http.Request) {
	nodes := make([]datasetData, 0)
	result := paging{}
	result.Nodes = make([]datasetData, 0)
	datasets := make([]model.Dataset, 0)

	offset, limit := paginationx.Parse(r.URL.Query())

	// Filters
	u, _ := url.Parse(r.URL.String())
	queryMap := u.Query()

	searchQuery := r.URL.Query().Get("q")
	sort := r.URL.Query().Get("sort")

	filters := util.GenerateFilters(queryMap["tag"])
	filteredDatasetIDs := make([]uint, 0)
	var err error

	if filters != "" || searchQuery != "" {
		// Search dataset with filter
		var hits []interface{}
		var res map[string]interface{}

		if searchQuery != "" {
			hits, err = meilisearchx.SearchWithQuery("mande", searchQuery, filters, "dataset")
		} else {
			res, err = meilisearchx.SearchWithoutQuery("mande", filters, "dataset")
			if _, found := res["hits"]; found {
				hits = res["hits"].([]interface{})
			}
		}
		if err != nil {
			loggerx.Error(err)
			renderx.JSON(w, http.StatusOK, result)
			return
		}

		filteredDatasetIDs = meilisearchx.GetIDArray(hits)
		if len(filteredDatasetIDs) == 0 {
			renderx.JSON(w, http.StatusOK, result)
			return
		}
	}
	if sort != "asc" {
		sort = "desc"
	}

	tx := model.DB.Model(&model.Dataset{}).Preload("FeaturedMedium").Preload("Currency").Preload("Tags").Order("created_at " + sort)

	if len(filteredDatasetIDs) > 0 {
		err = tx.Where(filteredDatasetIDs).Count(&result.Total).Offset(offset).Limit(limit).Find(&datasets).Error
	} else {
		err = tx.Count(&result.Total).Offset(offset).Limit(limit).Find(&datasets).Error
	}

	for _, dataset := range datasets {
		var formats []model.DatasetFormat

		data := &datasetData{}
		data.Formats = make([]model.DatasetFormat, 0)
		model.DB.Model(&model.DatasetFormat{}).Where(&model.DatasetFormat{
			DatasetID: uint(dataset.ID),
		}).Preload("Format").Find(&formats)

		data.Formats = append(data.Formats, formats...)

		data.Dataset = dataset

		nodes = append(nodes, *data)
	}

	result.Nodes = nodes

	renderx.JSON(w, http.StatusOK, result)
}
