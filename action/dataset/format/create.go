package format

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
	"github.com/go-chi/chi"
)

// create - Create dataset format
// @Summary Create dataset format
// @Description Create dataset format
// @Tags Dataset Format
// @ID add-dataset-format
// @Consume json
// @Produce  json
// @Param dataset_id path string true "Dataset ID"
// @Param DatasetFormat body datasetFormat true "Dataset Format object"
// @Success 201 {object} model.DatasetFormat
// @Failure 400 {array} string
// @Router /datasets/{dataset_id}/format [post]
func create(w http.ResponseWriter, r *http.Request) {

	datasetID := chi.URLParam(r, "dataset_id")
	id, err := strconv.Atoi(datasetID)

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	datasetFormat := &datasetFormat{}

	json.NewDecoder(r.Body).Decode(&datasetFormat)

	validationError := validationx.Check(datasetFormat)
	if validationError != nil {
		errorx.Render(w, validationError)
		return
	}

	result := &model.DatasetFormat{}

	result.FormatID = datasetFormat.FormatID
	result.DatasetID = uint(id)
	result.URL = datasetFormat.URL

	err = model.DB.Model(&model.DatasetFormat{}).Create(&result).Error

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	model.DB.Model(&model.DatasetFormat{}).Preload("Format").First(&result)

	renderx.JSON(w, http.StatusCreated, result)
}
