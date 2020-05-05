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
// @Param id path string true "Product ID"
// @Param tid path string true "ProductTag ID"
// @Success 200
// @Failure 400 {array} string
// @Router /products/{id}/tag/{tid} [delete]
func delete(w http.ResponseWriter, r *http.Request) {

	productTagID := chi.URLParam(r, "tid")
	tid, err := strconv.Atoi(productTagID)

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

	productTags := &model.ProductTag{}
	productTags.ID = uint(tid)
	productTags.ProductID = uint(pid)

	// check record exists or not
	err = model.DB.First(&productTags).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	model.DB.Delete(&productTags)

	render.JSON(w, http.StatusOK, nil)
}
