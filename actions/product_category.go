package actions

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/models"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

// ProductCategory request body
type productCategory struct {
	CategoryID uint `json:"category_id"`
}

// CreateProductCategory - create productCategory
// @Summary Create productCategory
// @Description create productCategory
// @Tags ProductCategory
// @ID add-productCategory
// @Consume json
// @Produce  json
// @Param id path string true "Product ID"
// @Param ProductCategory body productCategory true "ProductCategory object"
// @Success 200 {object} models.ProductCategory
// @Failure 400 {array} string
// @Router /products/{id}/productCategories [post]
func CreateProductCategory(w http.ResponseWriter, r *http.Request) {

	productID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(productID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &models.ProductCategory{
		ProductID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}
	err = models.DB.Model(&models.ProductCategory{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(req)
}

// DeleteProductCategory - Delete productCategory by id
// @Summary Delete a productCategory
// @Description Delete productCategory by ID
// @Tags ProductCategory
// @ID delete-productCategory-by-id
// @Consume  json
// @Param id path string true "Product ID"
// @Param cid path string true "ProductCategory ID"
// @Success 200 {object} models.ProductCategory
// @Failure 400 {array} string
// @Router /products/{id}/productCategories/{cid} [delete]
func DeleteProductCategory(w http.ResponseWriter, r *http.Request) {
	productCategoryID := chi.URLParam(r, "cid")
	id, err := strconv.Atoi(productCategoryID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	productCategory := &models.ProductCategory{
		ID: uint(id),
	}

	// check record exists or not
	err = models.DB.First(&productCategory).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	models.DB.Delete(&productCategory)

	json.NewEncoder(w).Encode(productCategory)
}
