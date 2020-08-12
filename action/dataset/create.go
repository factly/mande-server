package dataset

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
)

// create - create dataset
// @Summary Create dataset
// @Description create dataset
// @Tags Dataset
// @ID add-dataset
// @Consume json
// @Produce  json
// @Param Dataset body dataset true "Dataset object"
// @Success 201 {object} model.Dataset
// @Failure 400 {array} string
// @Router /datasets [post]
func create(w http.ResponseWriter, r *http.Request) {

	dataset := dataset{}

	json.NewDecoder(r.Body).Decode(&dataset)

	validationError := validationx.Check(dataset)
	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	result := &datasetData{}
	result.Tags = make([]model.Tag, 0)
	result.Dataset = model.Dataset{
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
	}

	err := model.DB.Model(&model.Dataset{}).Create(&result.Dataset).Error
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	model.DB.Preload("FeaturedMedium").First(&result.Dataset)

	// creating new dataset tags
	for _, id := range dataset.TagIDs {
		datasetTag := &model.DatasetTag{}
		datasetTag.TagID = uint(id)
		datasetTag.DatasetID = result.ID

		err = model.DB.Model(&model.DatasetTag{}).Create(&datasetTag).Error
		if err != nil {
			errorx.Render(w, errorx.Parser(errorx.DBError()))
			return
		}
	}

	// fetch all dataset tags
	datasetTags := []model.DatasetTag{}
	model.DB.Model(&model.DatasetTag{}).Where(&model.DatasetTag{
		DatasetID: result.Dataset.ID,
	}).Preload("Tag").Find(&datasetTags)

	// appending dataset tags to result
	for _, datasetTag := range datasetTags {
		result.Tags = append(result.Tags, datasetTag.Tag)
	}

	renderx.JSON(w, http.StatusCreated, result)
}
