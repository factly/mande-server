package actions

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"../models"
	"github.com/go-chi/chi"
)

func GetProduct(w http.ResponseWriter, r *http.Request) {
	productId := chi.URLParam(r, "productId")
	id, err := strconv.Atoi(productId)

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

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("create product")
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
	fmt.Printf("%+v", req.Status)
	json.NewEncoder(w).Encode(req)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productId := chi.URLParam(r, "productId")
	id, err := strconv.Atoi(productId)

	if err != nil {
		log.Fatal(err)
	}

	req := &models.Product{
		ID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&req)
	product := &models.Product{}
	models.DB.First(&models.Product{}).First(&req)

	if req.CurrencyID != 0 {
		product.CurrencyID = req.CurrencyID
	}
	if req.ProductTypeID != 0 {
		product.ProductTypeID = req.ProductTypeID
	}
	if req.StatusID != 0 {
		product.StatusID = req.StatusID
	}
	if req.Title != "" {
		product.Title = req.Title
	}
	if req.Price != 0 {
		product.Price = req.Price
	}
	if req.Slug != "" {
		product.Slug = req.Slug
	}

	models.DB.Model(&models.Product{}).Update(&product)

	models.DB.Model(&req).Association("ProductType").Find(&req.ProductType)
	models.DB.Model(&req).Association("Currency").Find(&req.Currency)
	models.DB.Model(&req).Association("Status").Find(&req.Status)

	json.NewEncoder(w).Encode(product)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productId := chi.URLParam(r, "productId")
	id, err := strconv.Atoi(productId)

	if err != nil {
		log.Fatal(err)
	}

	product := &models.Product{
		ID: uint(id),
	}
	json.NewDecoder(r.Body).Decode(&product)

	models.DB.First(&product)
	models.DB.Delete(&product)

	json.NewEncoder(w).Encode(product)
}
