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
	ProductID  uint `json:"product_id"`
}

// GetProductCategories - Get all productCategories
// @Summary Show all productCategories
// @Description Get all productCategories
// @Tags ProductCategory
// @ID get-all-productCategories
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} models.ProductCategory
// @Router /productCategories [get]
func GetProductCategories(w http.ResponseWriter, r *http.Request) {

	var productCategories []models.ProductCategory
	p := r.URL.Query().Get("page")
	pg, _ := strconv.Atoi(p) // pg contains page number
	l := r.URL.Query().Get("limit")
	li, _ := strconv.Atoi(l) // li contains perPage number

	offset := 0 // no. of records to skip
	limit := 5  // limt

	if li > 0 && li <= 10 {
		limit = li
	}

	if pg > 1 {
		offset = (pg - 1) * limit
	}

	models.DB.Offset(offset).Limit(limit).Model(&models.ProductCategory{}).Find(&productCategories)

	json.NewEncoder(w).Encode(productCategories)
}

// GetProductCategory - Get productCategory by id
// @Summary Show a productCategory by id
// @Description Get productCategory by ID
// @Tags ProductCategory
// @ID get-productCategory-by-id
// @Produce  json
// @Param id path string true "ProductCategory ID"
// @Success 200 {object} models.ProductCategory
// @Failure 400 {array} string
// @Router /productCategories/{id} [get]
func GetProductCategory(w http.ResponseWriter, r *http.Request) {

	productCategoryID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(productCategoryID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &models.ProductCategory{
		ID: uint(id),
	}

	err = models.DB.Model(&models.ProductCategory{}).First(&req).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(req)
}

// CreateProductCategory - create productCategory
// @Summary Create productCategory
// @Description create productCategory
// @Tags ProductCategory
// @ID add-productCategory
// @Consume json
// @Produce  json
// @Param ProductCategory body productCategory true "ProductCategory object"
// @Success 200 {object} models.ProductCategory
// @Failure 400 {array} string
// @Router /productCategories [post]
func CreateProductCategory(w http.ResponseWriter, r *http.Request) {

	req := &models.ProductCategory{}

	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.Struct(req)
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

// UpdateProductCategory - Update productCategory by id
// @Summary Update a productCategory by id
// @Description Update productCategory by ID
// @Tags ProductCategory
// @ID update-productCategory-by-id
// @Produce json
// @Consume json
// @Param id path string true "ProductCategory ID"
// @Param ProductCategory body productCategory false "ProductCategory"
// @Success 200 {object} models.ProductCategory
// @Failure 400 {array} string
// @Router /productCategories/{id} [put]
func UpdateProductCategory(w http.ResponseWriter, r *http.Request) {

	productCategoryID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(productCategoryID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &models.ProductCategory{}
	productCategory := &models.ProductCategory{
		ID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&req)

	models.DB.Model(&productCategory).Updates(models.ProductCategory{
		CategoryID: req.CategoryID,
		ProductID:  req.ProductID,
	})
	models.DB.First(&productCategory)

	json.NewEncoder(w).Encode(productCategory)
}

// DeleteProductCategory - Delete productCategory by id
// @Summary Delete a productCategory
// @Description Delete productCategory by ID
// @Tags ProductCategory
// @ID delete-productCategory-by-id
// @Consume  json
// @Param id path string true "ProductCategory ID"
// @Success 200 {object} models.ProductCategory
// @Failure 400 {array} string
// @Router /productCategories/{id} [delete]
func DeleteProductCategory(w http.ResponseWriter, r *http.Request) {

	productCategoryID := chi.URLParam(r, "id")
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
