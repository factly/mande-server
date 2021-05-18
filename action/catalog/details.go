package catalog

import (
	"net/http"
	"strconv"

	"github.com/factly/mande-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// details - Get catalog by id
// @Summary Show a catalog by id
// @Description Get catalog by ID
// @Tags Catalog
// @ID get-catalog-by-id
// @Produce  json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param catalog_id path string true "Catalog ID"
// @Success 200 {object} model.Catalog
// @Failure 400 {array} string
// @Router /catalogs/{catalog_id} [get]
func details(w http.ResponseWriter, r *http.Request) {
	result := &model.Catalog{}
	tx := model.DB.Model(&model.Catalog{}).Preload("FeaturedMedium").Preload("Products").Preload("Products.Currency").Preload("Products.FeaturedMedium").Preload("Products.Tags").Preload("Products.Datasets")

	catalogID := chi.URLParam(r, "catalog_id")
	id, err := strconv.Atoi(catalogID)

	if err != nil {
		tx.Where(&model.Catalog{
			Title: catalogID,
		})
	} else {
		tx.Where(&model.Catalog{
			Base: model.Base{
				ID: uint(id),
			},
		})
	}

	err = tx.First(&result).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	renderx.JSON(w, http.StatusOK, result)
}
