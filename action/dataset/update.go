package dataset

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/array"
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
	datasetTags := []model.DatasetTag{}

	json.NewDecoder(r.Body).Decode(&dataset)

	model.DB.Model(&result.Dataset).Updates(model.Dataset{
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
	}).Preload("FeaturedMedium").First(&result.Dataset)

	model.DB.Model(&model.DatasetFormat{}).Where(&model.DatasetFormat{
		DatasetID: uint(id),
	}).Preload("Format").Find(&result.Formats)

	// fetch existing dataset tags
	model.DB.Model(&model.DatasetTag{}).Where(&model.DatasetTag{
		DatasetID: uint(id),
	}).Preload("Tag").Find(&datasetTags)

	prevTagIDs := make([]uint, 0)
	datasetTagIDs := make([]uint, 0)
	// key as tag_id & value as dataset_tag
	mapperDatasetTag := map[uint]model.DatasetTag{}

	for _, datasetTag := range datasetTags {
		mapperDatasetTag[datasetTag.TagID] = datasetTag
		prevTagIDs = append(prevTagIDs, datasetTag.TagID)
	}

	toCreateIDs, toDeleteIDs := array.Difference(prevTagIDs, dataset.TagIDs)

	// map dataset tag ids
	for _, id := range toDeleteIDs {
		datasetTagIDs = append(datasetTagIDs, mapperDatasetTag[id].ID)
	}

	// delete dataset tags
	if len(datasetTagIDs) > 0 {
		model.DB.Where(datasetTagIDs).Delete(model.DatasetTag{})
	}

	// create dataset tags
	for _, id := range toCreateIDs {
		datasetTag := &model.DatasetTag{}
		datasetTag.TagID = uint(id)
		datasetTag.DatasetID = result.ID

		err = model.DB.Model(&model.DatasetTag{}).Create(&datasetTag).Error
		if err != nil {
			errorx.Render(w, errorx.Parser(errorx.DBError()))
			return
		}
	}

	// fetch updated dataset tags
	updatedDatasetTags := []model.DatasetTag{}
	model.DB.Model(&model.DatasetTag{}).Where(&model.DatasetTag{
		DatasetID: uint(id),
	}).Preload("Tag").Find(&updatedDatasetTags)

	// appending previous dataset tags to result
	for _, datasetTag := range updatedDatasetTags {
		result.Tags = append(result.Tags, datasetTag.Tag)
	}

	renderx.JSON(w, http.StatusOK, result)
}
