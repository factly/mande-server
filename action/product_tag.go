package action

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

// ProductTag request body
type productTags struct {
	TagID uint `json:"tag_id"`
}

// CreateProductTag - create productTags
// @Summary Create productTags
// @Description create productTags
// @Tags ProductTag
// @ID add-productTags
// @Consume json
// @Produce  json
// @Param id path string true "Product ID"
// @Param ProductTag body productTags true "ProductTag object"
// @Success 200 {object} model.ProductTag
// @Failure 400 {array} string
// @Router /products/{id}/tag [post]
func CreateProductTag(w http.ResponseWriter, r *http.Request) {

	productID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(productID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &model.ProductTag{
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

	err = model.DB.Model(&model.ProductTag{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(req)
}

// DeleteProductTag - Delete productTags by id
// @Summary Delete a productTags
// @Description Delete productTags by ID
// @Tags ProductTag
// @ID delete-productTags-by-id
// @Consume  json
// @Param id path string true "Product ID"
// @Param tid path string true "ProductTag ID"
// @Success 200 {object} model.ProductTag
// @Failure 400 {array} string
// @Router /products/{id}/tag/{tid} [delete]
func DeleteProductTag(w http.ResponseWriter, r *http.Request) {

	tagID := chi.URLParam(r, "tid")
	tid, err := strconv.Atoi(tagID)

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

	productTags := &model.ProductTag{
		TagID:     uint(tid),
		ProductID: uint(pid),
	}

	// check record exists or not
	err = model.DB.First(&productTags).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	model.DB.Delete(&productTags)

	json.NewEncoder(w).Encode(productTags)
}
