package actions

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/models"
	"github.com/go-chi/chi"
)

// currency request body
type currency struct {
	IsoCode string `json:"iso_code"`
	Name    string `json:"name"`
}

// GetCurrency - Get currency by id
// @Summary Show a currency by id
// @Description get currency by ID
// @Tags Currency
// @ID get-currency-by-id
// @Produce  json
// @Param id path string true "Currency ID"
// @Success 200 {object} models.Currency
// @Router /currencies/{id} [get]
func GetCurrency(w http.ResponseWriter, r *http.Request) {
	currencyID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(currencyID)
	if err != nil {
		log.Fatal(err)
	}

	req := &models.Currency{
		ID: uint(id),
	}

	models.DB.Model(&models.Currency{}).First(&req)

	json.NewEncoder(w).Encode(req)
}

// CreateCurrency - Create currency
// @Summary Create currency
// @Description Create currency
// @Tags Currency
// @ID add-currency
// @Consume json
// @Produce  json
// @Param Currency body currency true "Currency object"
// @Success 200 {object} models.Currency
// @Router /currencies [post]
func CreateCurrency(w http.ResponseWriter, r *http.Request) {

	req := &models.Currency{}

	json.NewDecoder(r.Body).Decode(&req)

	if validErrs := req.Validate(); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		w.Header().Set("Content-type", "applciation/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	err := models.DB.Model(&models.Currency{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(req)
}

// UpdateCurrency - Update currency by id
// @Summary Update a currency by id
// @Description Update currency by ID
// @Tags Currency
// @ID update-currency-by-id
// @Produce json
// @Consume json
// @Param id path string true "Currecny ID"
// @Param Currency body currency false "Currency"
// @Success 200 {object} models.Currency
// @Router /currencies/{id} [put]
func UpdateCurrency(w http.ResponseWriter, r *http.Request) {
	currencyID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(currencyID)
	if err != nil {
		log.Fatal(err)
	}

	currency := &models.Currency{
		ID: uint(id),
	}
	req := &models.Currency{}

	json.NewDecoder(r.Body).Decode(&req)
	models.DB.Model(&currency).Updates(models.Currency{IsoCode: req.IsoCode, Name: req.Name})
	models.DB.First(&currency)
	json.NewEncoder(w).Encode(currency)
}

// DeleteCurrency - Delete currency by id
// @Summary Delete a currency
// @Description Delete currency by ID
// @Tags Currency
// @ID delete-currency-by-id
// @Consume  json
// @Param id path string true "Currency ID"
// @Success 200 {object} models.Currency
// @Router /currencies/{id} [delete]
func DeleteCurrency(w http.ResponseWriter, r *http.Request) {
	currencyID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(currencyID)
	if err != nil {
		log.Fatal(err)
	}

	currency := &models.Currency{
		ID: uint(id),
	}

	models.DB.First(&currency)
	models.DB.Delete(&currency)

	json.NewEncoder(w).Encode(currency)
}
