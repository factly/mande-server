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

// Category request body
type category struct {
	Title    string `json:"title"`
	Slug     string `json:"slug"`
	ParentID uint   `json:"parent_id"`
}

// GetCategories - Get all categories
// @Summary Show all categories
// @Description Get all categories
// @Tags Category
// @ID get-all-categories
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} models.Category
// @Router /categories [get]
func GetCategories(w http.ResponseWriter, r *http.Request) {

	var categories []models.Category
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

	models.DB.Offset(offset).Limit(limit).Model(&models.Category{}).Find(&categories)

	json.NewEncoder(w).Encode(categories)
}

// GetCategory - Get category by id
// @Summary Show a category by id
// @Description Get category by ID
// @Tags Category
// @ID get-category-by-id
// @Produce  json
// @Param id path string true "Category ID"
// @Success 200 {object} models.Category
// @Failure 400 {array} string
// @Router /categories/{id} [get]
func GetCategory(w http.ResponseWriter, r *http.Request) {

	categoryID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(categoryID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &models.Category{
		ID: uint(id),
	}

	err = models.DB.Model(&models.Category{}).First(&req).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(req)
}

// CreateCategory - create category
// @Summary Create category
// @Description create category
// @Tags Category
// @ID add-category
// @Consume json
// @Produce  json
// @Param Category body category true "Category object"
// @Success 200 {object} models.Category
// @Failure 400 {array} string
// @Router /categories [post]
func CreateCategory(w http.ResponseWriter, r *http.Request) {

	req := &models.Category{}

	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = models.DB.Model(&models.Category{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(req)
}

// UpdateCategory - Update category by id
// @Summary Update a category by id
// @Description Update category by ID
// @Tags Category
// @ID update-category-by-id
// @Produce json
// @Consume json
// @Param id path string true "Category ID"
// @Param Category body category false "Category"
// @Success 200 {object} models.Category
// @Failure 400 {array} string
// @Router /categories/{id} [put]
func UpdateCategory(w http.ResponseWriter, r *http.Request) {

	categoryID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(categoryID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &models.Category{}
	category := &models.Category{
		ID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&req)

	models.DB.Model(&category).Updates(models.Category{
		Title: req.Title,
		Slug:  req.Slug,
	})
	models.DB.First(&category)

	json.NewEncoder(w).Encode(category)
}

// DeleteCategory - Delete category by id
// @Summary Delete a category
// @Description Delete category by ID
// @Tags Category
// @ID delete-category-by-id
// @Consume  json
// @Param id path string true "Category ID"
// @Success 200 {object} models.Category
// @Failure 400 {array} string
// @Router /categories/{id} [delete]
func DeleteCategory(w http.ResponseWriter, r *http.Request) {

	categoryID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(categoryID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	category := &models.Category{
		ID: uint(id),
	}

	// check record exists or not
	err = models.DB.First(&category).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	models.DB.Delete(&category)

	json.NewEncoder(w).Encode(category)
}
