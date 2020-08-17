package dataset

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
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

	validationError := validationx.Check(dataset)
	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	// check record exist or not
	err = model.DB.Preload("Tags").First(&result.Dataset).Error
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	tx := model.DB.Begin()

	oldTags := result.Tags
	newTags := make([]model.Tag, 0)
	model.DB.Model(&model.Tag{}).Where(dataset.TagIDs).Find(&newTags)

	if len(oldTags) > 0 {
		err = tx.Model(&result).Association("Tags").Delete(oldTags).Error
		if err != nil {
			tx.Rollback()
			loggerx.Error(err)
			errorx.Render(w, errorx.Parser(errorx.DBError()))
			return
		}
	}
	if len(newTags) == 0 {
		newTags = nil
	}

	if dataset.FeaturedMediumID == 0 {
		err = tx.Model(result.Dataset).Updates(map[string]interface{}{"featured_medium_id": nil}).First(&result.Dataset).Error
		result.FeaturedMediumID = 0
		if err != nil {
			tx.Rollback()
			loggerx.Error(err)
			errorx.Render(w, errorx.Parser(errorx.DBError()))
			return
		}
	}

	err = tx.Model(&result.Dataset).Set("gorm:association_autoupdate", false).Updates(model.Dataset{
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
		Price:            dataset.Price,
		CurrencyID:       dataset.CurrencyID,
		FeaturedMediumID: dataset.FeaturedMediumID,
		Tags:             newTags,
	}).Preload("FeaturedMedium").Preload("Tags").Preload("Currency").First(&result.Dataset).Error

	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	model.DB.Model(&model.DatasetFormat{}).Where(&model.DatasetFormat{
		DatasetID: uint(id),
	}).Preload("Format").Find(&result.Formats)

	tx.Commit()

	renderx.JSON(w, http.StatusOK, result)
}
