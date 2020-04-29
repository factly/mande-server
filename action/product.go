package action

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/model"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

type product struct {
	Title         string `json:"title"`
	Slug          string `json:"slug"`
	Price         int    `json:"price"`
	ProductTypeID uint   `json:"product_type_id"`
	StatusID      uint   `json:"status_id"`
	CurrencyID    uint   `json:"currency_id"`
}

// GetProducts - Get all products
// @Summary Show all products
// @Description Get all products
// @Tags Product
// @ID get-all-products
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} model.Product
// @Router /products [get]
func GetProducts(w http.ResponseWriter, r *http.Request) {

	var products []model.Product
	p := r.URL.Query().Get("page")
	pg, _ := strconv.Atoi(p) // pg contains page number
	l := r.URL.Query().Get("limit")
	li, _ := strconv.Atoi(l) // li contains perPage number

	offset := 0 // no. of records to skip
	limit := 5  // limt

	if li > 0 && li <= 10 {
		limit = li
	}

	if pg > 1 {
		offset = (pg - 1) * limit
	}

	model.DB.Offset(offset).Limit(limit).Preload("Currency").Preload("Status").Preload("ProductType").Model(&model.Product{}).Find(&products)

	json.NewEncoder(w).Encode(products)
}

// GetProduct - Get product by id
// @Summary Show a product by id
// @Description Get product by ID
// @Tags Product
// @ID get-product-by-id
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200 {object} model.Product
// @Failure 400 {array} string
// @Router /products/{id} [get]
func GetProduct(w http.ResponseWriter, r *http.Request) {

	productID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(productID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &model.Product{
		ID: uint(id),
	}

	err = model.DB.Model(&model.Product{}).First(&req).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	model.DB.Model(&req).Association("ProductType").Find(&req.ProductType)
	model.DB.Model(&req).Association("Currency").Find(&req.Currency)
	model.DB.Model(&req).Association("Status").Find(&req.Status)
	json.NewEncoder(w).Encode(req)
}

// CreateProduct - Create product
// @Summary Create product
// @Description Create product
// @Tags Product
// @ID add-product
// @Consume json
// @Produce  json
// @Param Product body product true "Product object"
// @Success 200 {object} model.Product
// @Failure 400 {array} string
// @Router /products [post]
func CreateProduct(w http.ResponseWriter, r *http.Request) {

	req := &model.Product{}
	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.StructExcept(req, "ProductType", "Status", "Currency")
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = model.DB.Model(&model.Product{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}
	model.DB.Model(&req).Association("ProductType").Find(&req.ProductType)
	model.DB.Model(&req).Association("Currency").Find(&req.Currency)
	model.DB.Model(&req).Association("Status").Find(&req.Status)
	json.NewEncoder(w).Encode(req)
}

// UpdateProduct - Update product by id
// @Summary Update a product by id
// @Description Update product by ID
// @Tags Product
// @ID update-product-by-id
// @Produce json
// @Consume json
// @Param id path string true "Product ID"
// @Param Product body product false "Product"
// @Success 200 {object} model.Product
// @Failure 400 {array} string
// @Router /products/{id} [put]
func UpdateProduct(w http.ResponseWriter, r *http.Request) {

	productID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(productID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	product := &model.Product{
		ID: uint(id),
	}

	req := &model.Product{}
	json.NewDecoder(r.Body).Decode(&req)

	model.DB.Model(&product).Updates(&model.Product{
		CurrencyID:    req.CurrencyID,
		ProductTypeID: req.ProductTypeID,
		StatusID:      req.StatusID,
		Title:         req.Title,
		Price:         req.Price,
		Slug:          req.Slug,
	})
	model.DB.First(&product).First(&product)
	model.DB.Model(&product).Association("ProductType").Find(&product.ProductType)
	model.DB.Model(&product).Association("Currency").Find(&product.Currency)
	model.DB.Model(&product).Association("Status").Find(&product.Status)

	json.NewEncoder(w).Encode(product)
}

// DeleteProduct - Delete product by id
// @Summary Delete a product
// @Description Delete product by ID
// @Tags Product
// @ID delete-product-by-id
// @Consume  json
// @Param id path string true "Product ID"
// @Success 200 {object} model.Product
// @Failure 400 {array} string
// @Router /products/{id} [delete]
func DeleteProduct(w http.ResponseWriter, r *http.Request) {

	productID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(productID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	product := &model.Product{
		ID: uint(id),
	}

	// check record exists or not
	err = model.DB.First(&product).Error
	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	model.DB.Model(&product).Association("ProductType").Find(&product.ProductType)
	model.DB.Model(&product).Association("Currency").Find(&product.Currency)
	model.DB.Model(&product).Association("Status").Find(&product.Status)
	model.DB.Delete(&product)

	json.NewEncoder(w).Encode(product)
}
