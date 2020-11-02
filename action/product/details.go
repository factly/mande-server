package product

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

type productRes struct {
	model.Product
	Memberships []model.Membership `json:"memberships"`
}

// userDetails - Get product by id
// @Summary Show a product by id
// @Description Get product by ID
// @Tags Product
// @ID get-product-by-id
// @Produce  json
// @Param X-User header string false "User ID"
// @Param product_id path string true "Product ID"
// @Success 200 {object} productRes
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

	result := &productRes{}
	result.Tags = make([]model.Tag, 0)
	result.Datasets = make([]model.Dataset, 0)
	result.Memberships = make([]model.Membership, 0)

	result.ID = uint(id)

	err = model.DB.Model(&model.Product{}).Preload("Currency").Preload("FeaturedMedium").Preload("Tags").Preload("Datasets").First(&result.Product).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	// Fetch all catalogs related to product
	catalogs := make([]model.Catalog, 0)
	_ = model.DB.Model(&result).Association("Catalogs").Find(&catalogs)

	catalogIDs := make([]uint, 0)
	for _, catalog := range catalogs {
		catalogIDs = append(catalogIDs, catalog.ID)
	}

	// Fetch all plans related to catalogs
	plans := make([]model.Plan, 0)
	model.DB.Model(&model.Plan{}).Joins("INNER JOIN dp_plan_catalog ON dp_plan_catalog.plan_id = dp_plan.id").Where("catalog_id IN (?)", catalogIDs).Find(&plans)

	planIDs := make([]uint, 0)
	for _, plan := range plans {
		planIDs = append(planIDs, plan.ID)
	}

	// Fetch memberships related to planIDs and user
	model.DB.Model(&model.Membership{}).Preload("Plan").Where("plan_id IN (?) AND user_id = ?", planIDs, uID).Find(&result.Memberships)

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
