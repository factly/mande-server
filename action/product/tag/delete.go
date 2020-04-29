package tag

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
)

// delete - Delete productTags by id
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
func delete(w http.ResponseWriter, r *http.Request) {

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
