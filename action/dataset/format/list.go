package format

import (
	"net/http"
	"strconv"

	"github.com/factly/mande-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// list response
type paging struct {
	Total int64                 `json:"total"`
	Nodes []model.DatasetFormat `json:"nodes"`
}

// list - Get all datsets format
// @Summary Show all datsets format
// @Description Get all datsets format
// @Tags Dataset Format
// @ID get-all-datsets-formats
// @Produce  json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param dataset_id path string true "Dataset ID"
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /datasets/{dataset_id}/format [get]
func list(w http.ResponseWriter, r *http.Request) {

	datasetID := chi.URLParam(r, "dataset_id")
	id, err := strconv.Atoi(datasetID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := paging{}
	result.Nodes = make([]model.DatasetFormat, 0)

	offset, limit := paginationx.Parse(r.URL.Query())

	// Check if logged in user owns the dataset

	model.DB.Preload("Format").Model(&model.DatasetFormat{}).Where(&model.DatasetFormat{
		DatasetID: uint(id),
	}).Count(&result.Total).Offset(offset).Limit(limit).Find(&result.Nodes)

	renderx.JSON(w, http.StatusOK, result)
}
