package category

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// delete - Delete productCategory by id
// @Summary Delete a productCategory
// @Description Delete productCategory by ID
// @Tags ProductCategory
// @ID delete-productCategory-by-id
// @Consume  json
// @Param product_id path string true "Product ID"
// @Param category_id path string true "Category ID"
// @Success 200
// @Failure 400 {array} string
// @Router /products/{product_id}/category/{category_id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {
	productCategoryID := chi.URLParam(r, "category_id")
	cid, err := strconv.Atoi(productCategoryID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	productID := chi.URLParam(r, "product_id")
	pid, err := strconv.Atoi(productID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	productCategory := &model.ProductCategory{}
	productCategory.ID = uint(cid)
	productCategory.ProductID = uint(pid)

	// check record exists or not
	err = model.DB.First(&productCategory).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	model.DB.Delete(&productCategory)

	render.JSON(w, http.StatusOK, nil)
}
