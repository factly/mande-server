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

// list - Get all datsets
// @Summary Show all datsets
// @Description Get all datsets
// @Tags Dataset
// @ID get-all-datsets
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /datasets [get]
func list(w http.ResponseWriter, r *http.Request) {
	nodes := []datasetData{}
	result := paging{}
	datasets := []model.Dataset{}

	offset, limit := paginationx.Parse(r.URL.Query())

	model.DB.Model(&model.Dataset{}).Count(&result.Total).Offset(offset).Limit(limit).Find(&datasets)

	for _, dataset := range datasets {
		var formats []model.DatasetFormat

		data := &datasetData{}
		model.DB.Model(&model.DatasetFormat{}).Where(&model.DatasetFormat{
			DatasetID: uint(dataset.ID),
		}).Preload("Format").Find(&formats)

		for _, f := range formats {
			data.Formats = append(data.Formats, f.Format)
		}

		data.Dataset = dataset

		nodes = append(nodes, *data)
	}

	result.Nodes = nodes

	renderx.JSON(w, http.StatusOK, result)
}
