package dataset

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/data-portal-server/util"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// userDetails - Get dataset by id
// @Summary Show a dataset by id
// @Description Get dataset by ID
// @Tags Dataset
// @ID get-dataset-by-id
// @Produce  json
// @Param X-User header string false "User ID"
// @Param dataset_id path string true "Dataset ID"
// @Success 200 {object} model.Dataset
// @Failure 400 {array} string
// @Router /datasets/{dataset_id} [get]
func userDetails(w http.ResponseWriter, r *http.Request) {

	uID, err := util.GetUser(r)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	datasetID := chi.URLParam(r, "dataset_id")
	id, err := strconv.Atoi(datasetID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &datasetData{}
	result.ID = uint(id)
	result.Formats = make([]model.DatasetFormat, 0)

	err = model.DB.Model(&model.Dataset{}).Preload("FeaturedMedium").Preload("Currency").Preload("Tags").First(&result.Dataset).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	// Check if the user owns dataset
	if checkOrderAssociation(uID, id) != 0 || checkMembershipAssociation(uID, id) != 0 {
		model.DB.Model(&model.DatasetFormat{}).Where(&model.DatasetFormat{
			DatasetID: uint(id),
		}).Preload("Format").Find(&result.Formats)
	}

	renderx.JSON(w, http.StatusOK, result)
}

// adminDetails - Get dataset by id
// @Summary Show a dataset by id
// @Description Get dataset by ID
// @Tags Dataset
// @ID get-dataset-by-id
// @Produce  json
// @Param X-User header string false "User ID"
// @Param dataset_id path string true "Dataset ID"
// @Success 200 {object} model.Dataset
// @Failure 400 {array} string
// @Router /datasets/{dataset_id} [get]
func adminDetails(w http.ResponseWriter, r *http.Request) {

	datasetID := chi.URLParam(r, "dataset_id")
	id, err := strconv.Atoi(datasetID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &datasetData{}
	result.ID = uint(id)
	result.Formats = make([]model.DatasetFormat, 0)

	err = model.DB.Model(&model.Dataset{}).Preload("FeaturedMedium").Preload("Currency").Preload("Tags").First(&result.Dataset).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	model.DB.Model(&model.DatasetFormat{}).Where(&model.DatasetFormat{
		DatasetID: uint(id),
	}).Preload("Format").Find(&result.Formats)

	renderx.JSON(w, http.StatusOK, result)
}

// check if the item is associated with any order of the user
func checkOrderAssociation(uID, id int) int {
	var count int

	orders := make([]model.Order, 0)
	model.DB.Model(&model.Order{}).Where(&model.Order{
		UserID: uint(uID),
	}).Find(&orders)

	orderIDs := make([]uint, 0)
	for _, order := range orders {
		orderIDs = append(orderIDs, order.ID)
	}

	products := make([]model.Product, 0)
	model.DB.Model(&model.Product{}).Joins("INNER JOIN dp_order_item ON dp_order_item.product_id = dp_product.id").Where("order_id IN (?)", orderIDs).Find(&products)

	productIDs := make([]uint, 0)
	for _, product := range products {
		productIDs = append(productIDs, product.ID)
	}

	dataset := model.Dataset{}
	dataset.ID = uint(id)

	count = model.DB.Model(&dataset).Where("product_id IN (?)", productIDs).Association("Products").Count()

	return count
}

// check if the item is associated with any membership of the user
func checkMembershipAssociation(uID, id int) int {
	var count int
	// Get all plans of user from membership
	memberships := make([]model.Membership, 0)

	model.DB.Model(&model.Membership{}).Where(&model.Membership{
		UserID: uint(uID),
	}).Find(&memberships)

	planIDs := make([]uint, 0)
	for _, membership := range memberships {
		planIDs = append(planIDs, membership.PlanID)
	}

	// Get all catalogs associated with all the plans
	catalogs := make([]model.Catalog, 0)
	model.DB.Model(&model.Catalog{}).Joins("INNER JOIN dp_plan_catalog ON dp_plan_catalog.catalog_id = dp_catalog.id").Where("plan_id IN (?)", planIDs).Find(&catalogs)

	catalogIDs := make([]uint, 0)
	for _, catalog := range catalogs {
		catalogIDs = append(catalogIDs, catalog.ID)
	}

	// Get all products associated with all the catalogs
	products := make([]model.Product, 0)
	model.DB.Model(&model.Product{}).Joins("INNER JOIN dp_catalog_product ON dp_catalog_product.product_id = dp_product.id").Where("catalog_id IN (?)", catalogIDs).Find(&products)

	productIDs := make([]uint, 0)
	for _, product := range products {
		productIDs = append(productIDs, product.ID)
	}

	// Count number of datasets in all products
	model.DB.Model(&model.Product{}).Joins("INNER JOIN dp_product_dataset ON dp_product_dataset.product_id = dp_product.id").Where("product_id IN (?) AND dataset_id = ?", productIDs, id).Count(&count)

	return count
}
