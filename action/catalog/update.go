package catalog

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// update - Update catalog by id
// @Summary Update a catalog by id
// @Description Update catalog by ID
// @Tags Catalog
// @ID update-catalog-by-id
// @Produce json
// @Consume json
// @Param catalog_id path string true "Catalog ID"
// @Param Catalog body catalog false "Catalog"
// @Success 200 {object} model.Catalog
// @Failure 400 {array} string
// @Router /catalogs/{catalog_id} [put]
func update(w http.ResponseWriter, r *http.Request) {

	catalogID := chi.URLParam(r, "catalog_id")
	id, err := strconv.Atoi(catalogID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	catalog := &catalog{}
	result := &model.Catalog{}
	result.ID = uint(id)
	result.Products = make([]model.Product, 0)

	json.NewDecoder(r.Body).Decode(&catalog)

	// check record exist or not
	err = model.DB.Model(&model.Catalog{}).Preload("FeaturedMedium").Preload("Products").First(&result).Error
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	oldProducts := result.Products
	newProducts := make([]model.Product, 0)
	model.DB.Model(&model.Product{}).Where(catalog.ProductIDs).Find(&newProducts)

	if len(oldProducts) > 0 {
		model.DB.Model(&result).Association("Products").Delete(oldProducts)
	}
	if len(newProducts) == 0 {
		newProducts = nil
	}

	model.DB.Model(&result).Set("gorm:association_autoupdate", false).Updates(model.Catalog{
		Title:            catalog.Title,
		Description:      catalog.Description,
		FeaturedMediumID: catalog.FeaturedMediumID,
		Price:            catalog.Price,
		PublishedDate:    catalog.PublishedDate,
		Products:         newProducts,
	}).Preload("FeaturedMedium").Preload("Products").Preload("Products.Currency").Preload("Products.Tags").Preload("Products.Datasets").First(&result)

	renderx.JSON(w, http.StatusOK, result)
}
