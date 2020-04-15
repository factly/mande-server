package actions

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-api/models"
)

func CreateProductType(w http.ResponseWriter, r *http.Request) {

	req := &models.ProductType{}
	json.NewDecoder(r.Body).Decode(&req)

	if validErrs := req.Validate(); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		w.Header().Set("Content-type", "applciation/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	err := models.DB.Model(&models.ProductType{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(req)
}

func UpdateProductType(w http.ResponseWriter, r *http.Request) {

	req := &models.ProductType{}

	json.NewDecoder(r.Body).Decode(&req)
	productType := &models.ProductType{}
	models.DB.First(&models.ProductType{})

	if req.Name != "" {
		productType.Name = req.Name
	}

	models.DB.Model(&models.ProductType{}).Update(&productType)
	json.NewEncoder(w).Encode(productType)
}

func DeleteProductType(w http.ResponseWriter, r *http.Request) {
	productType := &models.ProductType{}
	json.NewDecoder(r.Body).Decode(&productType)

	models.DB.First(&productType)
	models.DB.Delete(&productType)

	json.NewEncoder(w).Encode(productType)
}
