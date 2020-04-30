package category

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// delete - Delete productCategory by id
// @Summary Delete a productCategory
// @Description Delete productCategory by ID
// @Tags ProductCategory
// @ID delete-productCategory-by-id
// @Consume  json
// @Param id path string true "Product ID"
// @Param cid path string true "Category ID"
// @Success 200 {object} model.ProductCategory
// @Failure 400 {array} string
// @Router /products/{id}/category/{cid} [delete]
func delete(w http.ResponseWriter, r *http.Request) {
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

	productCategory := &model.ProductCategory{
		CategoryID: uint(cid),
		ProductID:  uint(pid),
	}

	// check record exists or not
	err = model.DB.First(&productCategory).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	model.DB.Delete(&productCategory)

	json.NewEncoder(w).Encode(productCategory)
}
