package catalog

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
)

// create - create catalog
// @Summary Create catalog
// @Description create catalog
// @Tags Catalog
// @ID add-catalog
// @Consume json
// @Produce  json
// @Param Catalog body catalog true "Catalog object"
// @Success 201 {object} model.Catalog
// @Failure 400 {array} string
// @Router /catalogs [post]
func create(w http.ResponseWriter, r *http.Request) {

	catalog := catalog{}
	result := model.Catalog{}
	result.Products = make([]model.Product, 0)

	json.NewDecoder(r.Body).Decode(&catalog)

	validationError := validationx.Check(catalog)
	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	result = model.Catalog{
		Title:            catalog.Title,
		Description:      catalog.Description,
		FeaturedMediumID: catalog.FeaturedMediumID,
		Price:            catalog.Price,
		PublishedDate:    catalog.PublishedDate,
	}

	model.DB.Model(&model.Product{}).Where(catalog.ProductIDs).Find(&result.Products)

	err := model.DB.Model(&model.Catalog{}).Create(&result).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	model.DB.Preload("FeaturedMedium").Preload("Products").Preload("Products.Currency").First(&result)

	renderx.JSON(w, http.StatusCreated, result)
}
