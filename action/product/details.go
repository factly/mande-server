package product

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/util"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// userDetails - Get product by id
// @Summary Show a product by id
// @Description Get product by ID
// @Tags Product
// @ID get-product-by-id
// @Produce  json
// @Param X-User header string false "User ID"
// @Param product_id path string true "Product ID"
// @Success 200 {object} model.Product
// @Failure 400 {array} string
// @Router /products/{product_id} [get]
func userDetails(w http.ResponseWriter, r *http.Request) {

	uID, err := util.GetUser(r)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	productID := chi.URLParam(r, "product_id")
	id, err := strconv.Atoi(productID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	// check if product is owned by user
	if checkOrderAssociation(uID, id) == 0 && checkMembershipAssociation(uID, id) == 0 {
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	result := &model.Product{}
	result.Tags = make([]model.Tag, 0)
	result.Datasets = make([]model.Dataset, 0)

	result.ID = uint(id)

	err = model.DB.Model(&model.Product{}).Preload("Currency").Preload("FeaturedMedium").Preload("Tags").Preload("Datasets").First(&result).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	renderx.JSON(w, http.StatusOK, result)
}

// adminDetails - Get product by id
// @Summary Show a product by id
// @Description Get product by ID
// @Tags Product
// @ID get-product-by-id
// @Produce  json
// @Param X-User header string false "User ID"
// @Param product_id path string true "Product ID"
// @Success 200 {object} model.Product
// @Failure 400 {array} string
// @Router /products/{product_id} [get]
func adminDetails(w http.ResponseWriter, r *http.Request) {

	productID := chi.URLParam(r, "product_id")
	id, err := strconv.Atoi(productID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &model.Product{}
	result.Tags = make([]model.Tag, 0)
	result.Datasets = make([]model.Dataset, 0)

	result.ID = uint(id)

	err = model.DB.Model(&model.Product{}).Preload("Currency").Preload("FeaturedMedium").Preload("Tags").Preload("Datasets").First(&result).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	renderx.JSON(w, http.StatusOK, result)
}

// check if the item is associated with any order of the user
func checkOrderAssociation(uID, id int) int {
	var count int
	model.DB.Model(&model.Order{}).Joins("INNER JOIN dp_order_item ON dp_order_item.order_id = dp_order.id").Where("user_id = ? AND product_id = ?", uID, id).Count(&count)

	return count
}

// check if the item is associated with any membership of the user
func checkMembershipAssociation(uID, id int) int {
	memberships := make([]model.Membership, 0)

	model.DB.Model(&model.Membership{}).Where(&model.Membership{
		UserID: uint(uID),
	}).Find(&memberships)

	planIDs := make([]uint, 0)
	for _, membership := range memberships {
		planIDs = append(planIDs, membership.PlanID)
	}

	catalogs := make([]model.Catalog, 0)
	model.DB.Model(&model.Catalog{}).Joins("INNER JOIN dp_plan_catalog ON dp_plan_catalog.catalog_id = dp_catalog.id").Where("plan_id IN (?)", planIDs).Find(&catalogs)

	catalogIDs := make([]uint, 0)
	for _, catalog := range catalogs {
		catalogIDs = append(catalogIDs, catalog.ID)
	}

	var count int
	model.DB.Table("dp_catalog_product").Where("catalog_id IN (?) AND product_id = ?", catalogIDs, id).Count(&count)

	return count
}
