package actions

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/models"
	"github.com/factly/data-portal-api/validationerrors"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

// ProductCategory request body
type productCategory struct {
	CategoryID uint `json:"category_id"`
	ProductID  uint `json:"product_id"`
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
	w.Header().Set("Content-Type", "application/json")
	productCategoryID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(productCategoryID)

	if err != nil {
		validationerrors.InvalidID(w, r)
		return
	}

	req := &models.ProductCategory{
		ID: uint(id),
	}

	models.DB.Model(&models.ProductCategory{}).First(&req)

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
	w.Header().Set("Content-Type", "application/json")
	req := &models.ProductCategory{}

	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		msg := err.Error()
		validationerrors.ValidErrors(w, r, msg)
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
	w.Header().Set("Content-Type", "application/json")
	productCategoryID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(productCategoryID)

	if err != nil {
		validationerrors.InvalidID(w, r)
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
	w.Header().Set("Content-Type", "application/json")
	productCategoryID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(productCategoryID)

	if err != nil {
		validationerrors.InvalidID(w, r)
		return
	}

	productCategory := &models.ProductCategory{
		ID: uint(id),
	}

	// check record exists or not
	err = models.DB.First(&productCategory).Error

	if err != nil {
		validationerrors.RecordNotFound(w, r)
		return
	}
	models.DB.Delete(&productCategory)

	json.NewEncoder(w).Encode(productCategory)
}
