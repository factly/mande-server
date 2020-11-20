package dataset

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/meili"
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
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param Dataset body dataset true "Dataset object"
// @Success 201 {object} model.Dataset
// @Failure 400 {array} string
// @Router /datasets [post]
func create(w http.ResponseWriter, r *http.Request) {

	dataset := dataset{}

	err := json.NewDecoder(r.Body).Decode(&dataset)
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

	featuredMediumID := &dataset.FeaturedMediumID
	if dataset.FeaturedMediumID == 0 {
		featuredMediumID = nil
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
		SampleURL:        dataset.SampleURL,
		RelatedArticles:  dataset.RelatedArticles,
		TimeSaved:        dataset.TimeSaved,
		FeaturedMediumID: featuredMediumID,
		Price:            dataset.Price,
		CurrencyID:       dataset.CurrencyID,
	}

	if len(dataset.TagIDs) > 0 {
		model.DB.Model(&model.Tag{}).Where(dataset.TagIDs).Find(&result.Tags)
	}

	tx := model.DB.Begin()
	err = tx.Model(&model.Dataset{}).Create(&result.Dataset).Error
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	tx.Preload("FeaturedMedium").Preload("Currency").Preload("Tags").First(&result.Dataset)

	// Insert into meili index
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
	}

	err = meili.AddDocument(meiliObj)
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	tx.Commit()
	renderx.JSON(w, http.StatusCreated, result)
}
