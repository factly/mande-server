package dataset

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/factly/mande-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/meilisearchx"
	"github.com/factly/x/middlewarex"
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
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param dataset_id path string true "Dataset ID"
// @Param Dataset body dataset false "Dataset"
// @Success 200 {object} model.Dataset
// @Failure 400 {array} string
// @Router /datasets/{dataset_id} [put]
func update(w http.ResponseWriter, r *http.Request) {
	uID, err := middlewarex.GetUser(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	datasetID := chi.URLParam(r, "dataset_id")
	id, err := strconv.Atoi(datasetID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	dataset := dataset{}
	result := &datasetData{}
	result.ID = uint(id)
	result.Formats = make([]model.DatasetFormat, 0)

	err = json.NewDecoder(r.Body).Decode(&dataset)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	validationError := validationx.Check(dataset)
	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	// check record exist or not
	err = model.DB.First(&result.Dataset).Error
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	tx := model.DB.Begin()

	newTags := make([]model.Tag, 0)
	if len(dataset.TagIDs) > 0 {
		model.DB.Model(&model.Tag{}).Where(dataset.TagIDs).Find(&newTags)
		err = tx.Model(&result.Dataset).Association("Tags").Replace(&newTags)
	} else {
		err = tx.Model(&result.Dataset).Association("Tags").Clear()
	}

	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	featuredMediumID := &dataset.FeaturedMediumID
	if dataset.FeaturedMediumID == 0 {
		err = tx.Omit("Tags").Model(result.Dataset).Updates(map[string]interface{}{"featured_medium_id": nil}).Error
		featuredMediumID = nil
		if err != nil {
			tx.Rollback()
			loggerx.Error(err)
			errorx.Render(w, errorx.Parser(errorx.DBError()))
			return
		}
	}

	err = tx.Model(&result.Dataset).Omit("Tags").Updates(model.Dataset{
		Base:             model.Base{UpdatedByID: uint(uID)},
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
		SampleURL:        dataset.SampleURL,
		RelatedArticles:  dataset.RelatedArticles,
		TimeSaved:        dataset.TimeSaved,
		Price:            dataset.Price,
		CurrencyID:       dataset.CurrencyID,
		FeaturedMediumID: featuredMediumID,
		ProfilingURL:     dataset.ProfilingURL,
		IsPublic:         dataset.IsPublic,
	}).Preload("FeaturedMedium").Preload("Currency").Preload("Tags").First(&result.Dataset).Error

	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	model.DB.Model(&model.DatasetFormat{}).Where(&model.DatasetFormat{
		DatasetID: uint(id),
	}).Preload("Format").Find(&result.Formats)

	// Update into meili index
	meiliObj := map[string]interface{}{
		"id":            result.ID,
		"kind":          "dataset",
		"title":         result.Title,
		"description":   result.Description,
		"source":        result.Source,
		"frequency":     result.Frequency,
		"granuality":    result.Granularity,
		"contact_name":  result.ContactName,
		"contact_email": result.ContactEmail,
		"license":       result.License,
		"data_standard": result.DataStandard,
		"price":         result.Price,
		"currency_id":   result.CurrencyID,
		"tag_IDs":       dataset.TagIDs,
		"profiling_url": dataset.ProfilingURL,
		"is_public":     dataset.IsPublic,
	}

	err = meilisearchx.UpdateDocument("mande", meiliObj)
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	tx.Commit()
	renderx.JSON(w, http.StatusOK, result)
}
