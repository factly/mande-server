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

// ProductTag request body
type productTag struct {
	TagID     uint `json:"tag_id"`
	ProductID uint `json:"product_id"`
}

// GetProductTag - Get productTag by id
// @Summary Show a productTag by id
// @Description Get productTag by ID
// @Tags ProductTag
// @ID get-productTag-by-id
// @Produce  json
// @Param id path string true "ProductTag ID"
// @Success 200 {object} models.ProductTag
// @Failure 400 {array} string
// @Router /productTags/{id} [get]
func GetProductTag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	productTagID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(productTagID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &models.ProductTag{
		ID: uint(id),
	}

	models.DB.Model(&models.ProductTag{}).First(&req)

	json.NewEncoder(w).Encode(req)
}

// CreateProductTag - create productTag
// @Summary Create productTag
// @Description create productTag
// @Tags ProductTag
// @ID add-productTag
// @Consume json
// @Produce  json
// @Param ProductTag body productTag true "ProductTag object"
// @Success 200 {object} models.ProductTag
// @Failure 400 {array} string
// @Router /productTags [post]
func CreateProductTag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	req := &models.ProductTag{}

	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = models.DB.Model(&models.ProductTag{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(req)
}

// UpdateProductTag - Update productTag by id
// @Summary Update a productTag by id
// @Description Update productTag by ID
// @Tags ProductTag
// @ID update-productTag-by-id
// @Produce json
// @Consume json
// @Param id path string true "ProductTag ID"
// @Param ProductTag body productTag false "ProductTag"
// @Success 200 {object} models.ProductTag
// @Failure 400 {array} string
// @Router /productTags/{id} [put]
func UpdateProductTag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	productTagID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(productTagID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &models.ProductTag{}
	productTag := &models.ProductTag{
		ID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&req)

	models.DB.Model(&productTag).Updates(models.ProductTag{
		TagID:     req.TagID,
		ProductID: req.ProductID,
	})
	models.DB.First(&productTag)

	json.NewEncoder(w).Encode(productTag)
}

// DeleteProductTag - Delete productTag by id
// @Summary Delete a productTag
// @Description Delete productTag by ID
// @Tags ProductTag
// @ID delete-productTag-by-id
// @Consume  json
// @Param id path string true "ProductTag ID"
// @Success 200 {object} models.ProductTag
// @Failure 400 {array} string
// @Router /productTags/{id} [delete]
func DeleteProductTag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	productTagID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(productTagID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	productTag := &models.ProductTag{
		ID: uint(id),
	}

	// check record exists or not
	err = models.DB.First(&productTag).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	models.DB.Delete(&productTag)

	json.NewEncoder(w).Encode(productTag)
}
