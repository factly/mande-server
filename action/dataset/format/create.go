package format

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/factly/mande-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/middlewarex"
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
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param dataset_id path string true "Dataset ID"
// @Param DatasetFormat body datasetFormat true "Dataset Format object"
// @Success 201 {object} model.DatasetFormat
// @Failure 400 {array} string
// @Router /datasets/{dataset_id}/format [post]
func create(w http.ResponseWriter, r *http.Request) {
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

	datasetFormat := &datasetFormat{}

	err = json.NewDecoder(r.Body).Decode(&datasetFormat)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	validationError := validationx.Check(datasetFormat)
	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	result := &model.DatasetFormat{}

	result.FormatID = datasetFormat.FormatID
	result.DatasetID = uint(id)
	result.URL = datasetFormat.URL

	err = model.DB.WithContext(context.WithValue(r.Context(), userContext, uID)).Model(&model.DatasetFormat{}).Create(&result).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	model.DB.Model(&model.DatasetFormat{}).Preload("Format").First(&result)

	renderx.JSON(w, http.StatusCreated, result)
}
