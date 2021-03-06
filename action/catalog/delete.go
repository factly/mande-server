package catalog

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/meilisearchx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// delete - Delete catalog by id
// @Summary Delete a catalog
// @Description Delete catalog by ID
// @Tags Catalog
// @ID delete-catalog-by-id
// @Consume  json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param catalog_id path string true "Catalog ID"
// @Success 200
// @Failure 400 {array} string
// @Router /catalogs/{catalog_id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {

	catalogID := chi.URLParam(r, "catalog_id")
	id, err := strconv.Atoi(catalogID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &model.Catalog{}
	result.ID = uint(id)

	// check record exists or not
	err = model.DB.Preload("Products").First(&result).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	tx := model.DB.Begin()

	_ = tx.Model(&result).Association("Products").Delete(result.Products)

	err = tx.Delete(&result).Error
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	err = meilisearchx.DeleteDocument("data-portal", result.ID, "catalog")
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	tx.Commit()

	renderx.JSON(w, http.StatusOK, nil)
}
