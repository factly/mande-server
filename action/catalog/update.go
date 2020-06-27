package catalog

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
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

	products := []model.CatalogProduct{}

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	catalog := &catalog{}
	result := &catalogData{}
	result.ID = uint(id)
	result.Products = make([]model.Product, 0)

	json.NewDecoder(r.Body).Decode(&catalog)

	model.DB.Model(&result.Catalog).Updates(model.Catalog{
		Title:           catalog.Title,
		Description:     catalog.Description,
		FeaturedMediaID: catalog.FeaturedMediaID,
		Price:           catalog.Price,
		PublishedDate:   catalog.PublishedDate,
	}).Preload("FeaturedMedia").First(&result.Catalog)

	// fetch all products
	model.DB.Model(&model.CatalogProduct{}).Where(&model.CatalogProduct{
		CatalogID: uint(id),
	}).Preload("Product").Find(&products)

	// delete products
	for _, p := range products {
		present := false
		for _, id := range catalog.ProductIDs {
			if p.ProductID == id {
				present = true
			}
		}
		if present == false {
			model.DB.Where(&model.CatalogProduct{
				ProductID: p.ProductID,
				CatalogID: uint(id),
			}).Delete(model.CatalogProduct{})
		}
	}

	// creating new products
	for _, id := range catalog.ProductIDs {
		present := false
		for _, p := range products {
			if p.ProductID == id {
				present = true
				result.Products = append(result.Products, p.Product)
			}
		}
		if present == false {
			catalogProduct := &model.CatalogProduct{}
			catalogProduct.ProductID = uint(id)
			catalogProduct.ProductID = result.ID

			err = model.DB.Model(&model.CatalogProduct{}).Create(&catalogProduct).Error

			if err != nil {
				return
			}
			model.DB.Model(&model.CatalogProduct{}).Preload("Product").First(&catalogProduct)
			result.Products = append(result.Products, catalogProduct.Product)
		}
	}

	renderx.JSON(w, http.StatusOK, result)
}
