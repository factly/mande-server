package product

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util/meili"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// delete - Delete product by id
// @Summary Delete a product
// @Description Delete product by ID
// @Tags Product
// @ID delete-product-by-id
// @Consume  json
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Param product_id path string true "Product ID"
// @Success 200
// @Failure 400 {array} string
// @Router /products/{product_id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {

	productID := chi.URLParam(r, "product_id")
	id, err := strconv.Atoi(productID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &model.Product{}

	result.ID = uint(id)

	// check record exists or not
	err = model.DB.Preload("Tags").Preload("Datasets").First(&result).Error
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	var totAssociated int64
	product := new(model.Product)
	product.ID = uint(id)

	// check if product is associated with catalog
	totAssociated = model.DB.Model(product).Association("Catalogs").Count()

	if totAssociated != 0 {
		loggerx.Error(errors.New("product is associated with catalog"))
		errorx.Render(w, errorx.Parser(errorx.CannotSaveChanges()))
		return
	}

	// check if product is associated with order
	totAssociated = model.DB.Model(product).Association("Orders").Count()

	if totAssociated != 0 {
		loggerx.Error(errors.New("product is associated with order"))
		errorx.Render(w, errorx.Parser(errorx.CannotSaveChanges()))
		return
	}

	tx := model.DB.Begin()
	if len(result.Tags) > 0 {
		_ = tx.Model(&result).Association("Tags").Delete(result.Tags)
	}

	if len(result.Datasets) > 0 {
		_ = tx.Model(&result).Association("Datasets").Delete(result.Datasets)
	}

	err = tx.Delete(&result).Error
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	err = meili.DeleteDocument(result.ID, "product")
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	tx.Commit()

	renderx.JSON(w, http.StatusOK, nil)
}
