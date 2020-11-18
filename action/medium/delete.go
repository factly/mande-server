package medium

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/meili"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// delete - Delete medium by id
// @Summary Delete a medium
// @Description Delete medium by ID
// @Tags Medium
// @ID delete-medium-by-id
// @Consume  json
// @Param medium_id path string true "Medium ID"
// @Success 200
// @Failure 400 {array} string
// @Router /media/{medium_id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {

	mediumID := chi.URLParam(r, "medium_id")
	id, err := strconv.Atoi(mediumID)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &model.Medium{}
	result.ID = uint(id)

	// check record exists or not
	err = model.DB.First(&result).Error
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	uintID := uint(id)

	// check if medium is associated with catalog
	var totAssociated int64
	model.DB.Model(&model.Catalog{}).Where(&model.Catalog{
		FeaturedMediumID: &uintID,
	}).Count(&totAssociated)

	if totAssociated != 0 {
		loggerx.Error(errors.New("medium is associated with catalog"))
		errorx.Render(w, errorx.Parser(errorx.CannotSaveChanges()))
		return
	}

	// check if medium is associated with dataset
	model.DB.Model(&model.Dataset{}).Where(&model.Dataset{
		FeaturedMediumID: &uintID,
	}).Count(&totAssociated)

	if totAssociated != 0 {
		loggerx.Error(errors.New("medium is associated with dataset"))
		errorx.Render(w, errorx.Parser(errorx.CannotSaveChanges()))
		return
	}

	// check if medium is associated with product
	model.DB.Model(&model.Product{}).Where(&model.Product{
		FeaturedMediumID: &uintID,
	}).Count(&totAssociated)

	if totAssociated != 0 {
		loggerx.Error(errors.New("medium is associated with product"))
		errorx.Render(w, errorx.Parser(errorx.CannotSaveChanges()))
		return
	}

	tx := model.DB.Begin()
	tx.Delete(&result)

	err = meili.DeleteDocument(result.ID, "medium")
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	tx.Commit()
	renderx.JSON(w, http.StatusOK, nil)
}
