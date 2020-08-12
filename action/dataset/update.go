package dataset

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// update - Update dataset by id
// @Summary Update a dataset by id
// @Description Update dataset by ID
// @Tags Dataset
// @ID update-dataset-by-id
// @Produce json
// @Consume json
// @Param dataset_id path string true "Dataset ID"
// @Param Dataset body dataset false "Dataset"
// @Success 200 {object} model.Dataset
// @Failure 400 {array} string
// @Router /datasets/{dataset_id} [put]
func update(w http.ResponseWriter, r *http.Request) {

	datasetID := chi.URLParam(r, "dataset_id")
	id, err := strconv.Atoi(datasetID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	dataset := &dataset{}
	result := &datasetData{}
	result.ID = uint(id)
	result.Formats = make([]model.DatasetFormat, 0)

	json.NewDecoder(r.Body).Decode(&dataset)

	// check record exist or not
	err = model.DB.Model(&model.Dataset{}).Preload("FeaturedMedium").Preload("Tags").First(&result.Dataset).Error
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	oldTags := result.Tags
	newTags := make([]model.Tag, 0)
	model.DB.Model(&model.Tag{}).Where(dataset.TagIDs).Find(&newTags)

	if len(oldTags) > 0 {
		model.DB.Model(&result).Association("Tags").Delete(oldTags)
	}
	if len(newTags) == 0 {
		newTags = nil
	}

	model.DB.Model(&result.Dataset).Set("gorm:association_autoupdate", false).Updates(model.Dataset{
		Title:            dataset.Title,
		Description:      dataset.Description,
		Source:           dataset.Source,
		Frequency:        dataset.Frequency,
		TemporalCoverage: dataset.TemporalCoverage,
		Granularity:      dataset.Granularity,
		ContactName:      dataset.ContactName,
		ContactEmail:     dataset.ContactEmail,
		License:          dataset.License,
		DataStandard:     dataset.DataStandard,
		RelatedArticles:  dataset.RelatedArticles,
		TimeSaved:        dataset.TimeSaved,
		FeaturedMediumID: dataset.FeaturedMediumID,
		Tags:             newTags,
	}).Preload("FeaturedMedium").Preload("Tags").First(&result.Dataset)

	model.DB.Model(&model.DatasetFormat{}).Where(&model.DatasetFormat{
		DatasetID: uint(id),
	}).Preload("Format").Find(&result.Formats)

	renderx.JSON(w, http.StatusOK, result)
}
