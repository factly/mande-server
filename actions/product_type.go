package actions

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/factly/data-portal-api/models"
)

// productType request body
type productType struct {
	Name string `json:"name"`
}

// CreateProductType - Create product type
// @Summary Create product type
// @Description Create product type
// @Tags Type
// @ID add-type
// @Consume json
// @Produce  json
// @Param Type body productType true "Type object"
// @Success 200 {object} models.ProductType
// @Router /products/{id}/type [post]
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
