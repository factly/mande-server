package tag

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/render"
	"github.com/factly/data-portal-server/validation"
	"github.com/go-chi/chi"
)

// delete - Delete productTags by id
// @Summary Delete a productTags
// @Description Delete productTags by ID
// @Tags ProductTag
// @ID delete-productTags-by-id
// @Consume  json
// @Param product_id path string true "Product ID"
// @Param tag_id path string true "ProductTag ID"
// @Success 200
// @Failure 400 {array} string
// @Router /products/{product_id}/tag/{tag_id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {

	productTagID := chi.URLParam(r, "tag_id")
	tid, err := strconv.Atoi(productTagID)

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

	result := &model.ProductTag{}
	result.ID = uint(tid)
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
