package catalog

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/meilisearchx"
	"github.com/factly/x/middlewarex"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
	"github.com/go-chi/chi"
)

// update - Update catalog by id
// @Summary Update a catalog by id
// @Description Update catalog by ID
// @Tags Catalog
// @ID update-catalog-by-id
// @Produce json
// @Consume json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param catalog_id path string true "Catalog ID"
// @Param Catalog body catalog false "Catalog"
// @Success 200 {object} model.Catalog
// @Failure 400 {array} string
// @Router /catalogs/{catalog_id} [put]
func update(w http.ResponseWriter, r *http.Request) {
	uID, err := middlewarex.GetUser(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	catalogID := chi.URLParam(r, "catalog_id")
	id, err := strconv.Atoi(catalogID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	catalog := &catalog{}
	result := model.Catalog{}
	result.ID = uint(id)
	result.Products = make([]model.Product, 0)

	err = json.NewDecoder(r.Body).Decode(&catalog)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	validationError := validationx.Check(catalog)
	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	// check record exist or not
	err = model.DB.First(&result).Error
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	tx := model.DB.Begin()

	newProducts := make([]model.Product, 0)
	if len(catalog.ProductIDs) > 0 {
		model.DB.Model(&model.Product{}).Where(catalog.ProductIDs).Find(&newProducts)
		err = tx.Model(&result).Association("Products").Replace(&newProducts)
	} else {
		err = tx.Model(&result).Association("Products").Clear()
	}

	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	featuredMediumID := &catalog.FeaturedMediumID
	if catalog.FeaturedMediumID == 0 {
		err = tx.Omit("Products").Model(result).Updates(map[string]interface{}{"featured_medium_id": nil}).Error
		featuredMediumID = nil
		if err != nil {
			tx.Rollback()
			loggerx.Error(err)
			errorx.Render(w, errorx.Parser(errorx.DBError()))
			return
		}
	}

	err = tx.Omit("Products").Model(&result).Updates(model.Catalog{
		Base:             model.Base{UpdatedByID: uint(uID)},
		Title:            catalog.Title,
		Description:      catalog.Description,
		FeaturedMediumID: featuredMediumID,
		PublishedDate:    catalog.PublishedDate,
		Products:         newProducts,
	}).Preload("FeaturedMedium").Preload("Products").Preload("Products.Currency").Preload("Products.FeaturedMedium").Preload("Products.Tags").Preload("Products.Datasets").First(&result).Error

	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	// Update into meili index
	meiliObj := map[string]interface{}{
		"id":             result.ID,
		"kind":           "catalog",
		"title":          result.Title,
		"description":    result.Description,
		"published_date": result.PublishedDate.Unix(),
		"product_ids":    catalog.ProductIDs,
	}

	err = meilisearchx.UpdateDocument("data-portal", meiliObj)
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	tx.Commit()

	renderx.JSON(w, http.StatusOK, result)
}
