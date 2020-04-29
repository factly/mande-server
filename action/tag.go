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

// tag request body
type tag struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

// GetTags - Get all tags
// @Summary Show all tags
// @Description Get all tags
// @Tags Tag
// @ID get-all-tags
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} models.Tag
// @Router /tags [get]
func GetTags(w http.ResponseWriter, r *http.Request) {

	var tags []models.Tag
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

	models.DB.Offset(offset).Limit(limit).Model(&models.Tag{}).Find(&tags)

	json.NewEncoder(w).Encode(tags)
}

// GetTag - Get tag by id
// @Summary Show a tag by id
// @Description Get tag by ID
// @Tags Tag
// @ID get-tag-by-id
// @Produce  json
// @Param id path string true "Tag ID"
// @Success 200 {object} models.Tag
// @Failure 400 {array} string
// @Router /tags/{id} [get]
func GetTag(w http.ResponseWriter, r *http.Request) {

	tagID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(tagID)
	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &models.Tag{
		ID: uint(id),
	}

	err = models.DB.Model(&models.Tag{}).First(&req).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(req)
}

// CreateTag - Create tag
// @Summary Create tag
// @Description Create tag
// @Tags Tag
// @ID add-tag
// @Consume json
// @Produce  json
// @Param Tag body tag true "Tag object"
// @Success 200 {object} models.Tag
// @Failure 400 {array} string
// @Router /tags [post]
func CreateTag(w http.ResponseWriter, r *http.Request) {

	req := &models.Tag{}

	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}
	err = models.DB.Model(&models.Tag{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(req)
}

// UpdateTag - Update tag by id
// @Summary Update a tag by id
// @Description Update tag by ID
// @Tags Tag
// @ID update-tag-by-id
// @Produce json
// @Consume json
// @Param id path string true "Tag ID"
// @Param Tag body tag false "Tag"
// @Success 200 {object} models.Tag
// @Failure 400 {array} string
// @Router /tags/{id} [put]
func UpdateTag(w http.ResponseWriter, r *http.Request) {

	tagID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(tagID)
	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &models.Tag{}
	tag := &models.Tag{
		ID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&req)

	models.DB.Model(&tag).Update(&models.Tag{Title: req.Title, Slug: req.Slug})
	models.DB.First(&tag)
	json.NewEncoder(w).Encode(tag)
}

// DeleteTag - Delete tag by id
// @Summary Delete a tag
// @Description Delete tag by ID
// @Tags Tag
// @ID delete-tag-by-id
// @Consume  json
// @Param id path string true "Tag ID"
// @Success 200 {object} models.Tag
// @Failure 400 {array} string
// @Router /tags/{id} [delete]
func DeleteTag(w http.ResponseWriter, r *http.Request) {

	tagID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(tagID)
	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	tag := &models.Tag{
		ID: uint(id),
	}

	// check record exists or not
	err = models.DB.First(&tag).Error
	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	models.DB.Delete(&tag)

	json.NewEncoder(w).Encode(tag)
}
