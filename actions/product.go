package actions

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/models"
	"github.com/go-chi/chi"
)

type product struct {
	Title         string `json:"title"`
	Slug          string `json:"slug"`
	Price         int    `json:"price"`
	ProductTypeID uint   `json:"product_type_id"`
	StatusID      uint   `json:"status_id"`
	CurrencyID    uint   `json:"currency_id"`
}

// GetProduct - Get product by id
// @Summary Show a product by id
// @Description Get product by ID
// @Tags Product
// @ID get-product-by-id
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product
// @Router /products/{id} [get]
func GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	productID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(productID)

	if err != nil {
		log.Fatal(err)
	}

	req := &models.Product{
		ID: uint(id),
	}

	models.DB.Model(&models.Product{}).First(&req)

	models.DB.Model(&req).Association("ProductType").Find(&req.ProductType)
	models.DB.Model(&req).Association("Currency").Find(&req.Currency)
	models.DB.Model(&req).Association("Status").Find(&req.Status)
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
// @Success 200 {object} models.Product
// @Router /products [post]
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	req := &models.Product{}
	json.NewDecoder(r.Body).Decode(&req)

	if validErrs := req.Validate(); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		w.Header().Set("Content-type", "applciation/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	err := models.DB.Model(&models.Product{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}
	models.DB.Model(&req).Association("ProductType").Find(&req.ProductType)
	models.DB.Model(&req).Association("Currency").Find(&req.Currency)
	models.DB.Model(&req).Association("Status").Find(&req.Status)
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
// @Success 200 {object} models.Product
// @Router /products/{id} [put]
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	productID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(productID)

	if err != nil {
		log.Fatal(err)
	}

	product := &models.Product{
		ID: uint(id),
	}

	req := &models.Product{}
	json.NewDecoder(r.Body).Decode(&req)

	models.DB.Model(&product).Updates(&models.Product{
		CurrencyID:    req.CurrencyID,
		ProductTypeID: req.ProductTypeID,
		StatusID:      req.StatusID,
		Title:         req.Title,
		Price:         req.Price,
		Slug:          req.Slug,
	})
	models.DB.First(&product).First(&product)
	models.DB.Model(&product).Association("ProductType").Find(&product.ProductType)
	models.DB.Model(&product).Association("Currency").Find(&product.Currency)
	models.DB.Model(&product).Association("Status").Find(&product.Status)

	json.NewEncoder(w).Encode(product)
}

// DeleteProduct - Delete product by id
// @Summary Delete a product
// @Description Delete product by ID
// @Tags Product
// @ID delete-product-by-id
// @Consume  json
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product
// @Router /products/{id} [delete]
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	productID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(productID)

	if err != nil {
		log.Fatal(err)
	}

	product := &models.Product{
		ID: uint(id),
	}

	models.DB.First(&product)
	models.DB.Model(&product).Association("ProductType").Find(&product.ProductType)
	models.DB.Model(&product).Association("Currency").Find(&product.Currency)
	models.DB.Model(&product).Association("Status").Find(&product.Status)
	models.DB.Delete(&product)

	json.NewEncoder(w).Encode(product)
}
