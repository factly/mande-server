package dataset

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// delete - Delete dataset by id
// @Summary Delete a dataset
// @Description Delete dataset by ID
// @Tags Dataset
// @ID delete-dataset-by-id
// @Consume  json
// @Param dataset_id path string true "Dataset ID"
// @Success 200
// @Failure 400 {array} string
// @Router /datasets/{dataset_id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {

	datasetID := chi.URLParam(r, "dataset_id")
	id, err := strconv.Atoi(datasetID)

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &model.Dataset{}
	result.ID = uint(id)

	// check record exists or not
	err = model.DB.First(&result).Error

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}
	model.DB.Delete(&result)

	renderx.JSON(w, http.StatusOK, nil)
}
