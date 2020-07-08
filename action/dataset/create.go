package dataset

import (
	"encoding/json"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
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
		errorx.Render(w, validationError)
		return
	}

	result := &model.Dataset{
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

	err := model.DB.Model(&model.Dataset{}).Create(&result).Error

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	model.DB.Preload("FeaturedMedium").First(&result)

	renderx.JSON(w, http.StatusCreated, result)
}
