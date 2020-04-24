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
type productTags struct {
	TagID     uint `json:"tag_id"`
	ProductID uint `json:"product_id"`
}

// GetProductTags - Get all productTags
// @Summary Show all productTags
// @Description Get all productTags
// @Tags ProductTag
// @ID get-all-productTags
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} models.ProductTag
// @Router /productTags [get]
func GetProductTags(w http.ResponseWriter, r *http.Request) {

	var productTags []models.ProductTag
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

	models.DB.Offset(offset).Limit(limit).Model(&models.ProductTag{}).Find(&productTags)

	json.NewEncoder(w).Encode(productTags)
}

// GetProductTag - Get productTags by id
// @Summary Show a productTags by id
// @Description Get productTags by ID
// @Tags ProductTag
// @ID get-productTags-by-id
// @Produce  json
// @Param id path string true "ProductTag ID"
// @Success 200 {object} models.ProductTag
// @Failure 400 {array} string
// @Router /productTags/{id} [get]
func GetProductTag(w http.ResponseWriter, r *http.Request) {

	productTagID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(productTagID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &models.ProductTag{
		ID: uint(id),
	}

	err = models.DB.Model(&models.ProductTag{}).First(&req).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(req)
}

// CreateProductTag - create productTags
// @Summary Create productTags
// @Description create productTags
// @Tags ProductTag
// @ID add-productTags
// @Consume json
// @Produce  json
// @Param ProductTag body productTags true "ProductTag object"
// @Success 200 {object} models.ProductTag
// @Failure 400 {array} string
// @Router /productTags [post]
func CreateProductTag(w http.ResponseWriter, r *http.Request) {

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

// UpdateProductTag - Update productTags by id
// @Summary Update a productTags by id
// @Description Update productTags by ID
// @Tags ProductTag
// @ID update-productTags-by-id
// @Produce json
// @Consume json
// @Param id path string true "ProductTag ID"
// @Param ProductTag body productTags false "ProductTag"
// @Success 200 {object} models.ProductTag
// @Failure 400 {array} string
// @Router /productTags/{id} [put]
func UpdateProductTag(w http.ResponseWriter, r *http.Request) {

	productTagID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(productTagID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &models.ProductTag{}
	productTags := &models.ProductTag{
		ID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&req)

	models.DB.Model(&productTags).Updates(models.ProductTag{
		TagID:     req.TagID,
		ProductID: req.ProductID,
	})
	models.DB.First(&productTags)

	json.NewEncoder(w).Encode(productTags)
}

// DeleteProductTag - Delete productTags by id
// @Summary Delete a productTags
// @Description Delete productTags by ID
// @Tags ProductTag
// @ID delete-productTags-by-id
// @Consume  json
// @Param id path string true "ProductTag ID"
// @Success 200 {object} models.ProductTag
// @Failure 400 {array} string
// @Router /productTags/{id} [delete]
func DeleteProductTag(w http.ResponseWriter, r *http.Request) {

	productTagID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(productTagID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	productTags := &models.ProductTag{
		ID: uint(id),
	}

	// check record exists or not
	err = models.DB.First(&productTags).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	models.DB.Delete(&productTags)

	json.NewEncoder(w).Encode(productTags)
}
