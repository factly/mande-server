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
// @Param id path string true "Product ID"
// @Param cid path string true "Category ID"
// @Success 200
// @Failure 400 {array} string
// @Router /products/{id}/category/{cid} [delete]
func delete(w http.ResponseWriter, r *http.Request) {
	productCategoryID := chi.URLParam(r, "cid")
	cid, err := strconv.Atoi(productCategoryID)

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

	result := &model.ProductCategory{}
	result.ID = uint(cid)
	result.ProductID = uint(pid)

	// check record exists or not
	err = model.DB.First(&result).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	model.DB.Delete(&result)

	render.JSON(w, http.StatusOK, nil)
}
