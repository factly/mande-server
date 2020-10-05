package dataset

import (
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
)

// list response
type paging struct {
	Total int           `json:"total"`
	Nodes []datasetData `json:"nodes"`
}

// userlist - Get all datsets
// @Summary Show all datsets
// @Description Get all datsets
// @Tags Dataset
// @ID get-all-datsets
// @Produce  json
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

	model.DB.Preload("FeaturedMedium").Preload("Currency").Preload("Tags").Model(&model.Dataset{}).Count(&result.Total).Offset(offset).Limit(limit).Find(&datasets)

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
