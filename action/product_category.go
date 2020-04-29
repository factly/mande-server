package action

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
// @Router /products/{id}/category [post]
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
// @Param cid path string true "Category ID"
// @Success 200 {object} models.ProductCategory
// @Failure 400 {array} string
// @Router /products/{id}/category/{cid} [delete]
func DeleteProductCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "cid")
	cid, err := strconv.Atoi(categoryID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	productID := chi.URLParam(r, "id")
	pid, err := strconv.Atoi(productID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	productCategory := &models.ProductCategory{
		CategoryID: uint(cid),
		ProductID:  uint(pid),
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
