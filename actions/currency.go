package actions

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/models"
	"github.com/factly/data-portal-api/validation"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

// currency request body
type currency struct {
	IsoCode string `json:"iso_code"`
	Name    string `json:"name"`
}

// GetCurrencies - Get all currencies
// @Summary Show all currencies
// @Description Get all currencies
// @Tags Currency
// @ID get-all-currencies
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} models.Currency
// @Router /currencies [get]
func GetCurrencies(w http.ResponseWriter, r *http.Request) {

	var currencies []models.Currency
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

	models.DB.Offset(offset).Limit(limit).Model(&models.Currency{}).Find(&currencies)

	json.NewEncoder(w).Encode(currencies)
}

// GetCurrency - Get currency by id
// @Summary Show a currency by id
// @Description get currency by ID
// @Tags Currency
// @ID get-currency-by-id
// @Produce  json
// @Param id path string false "Currency ID"
// @Success 200 {object} models.Currency
// @Failure 400 {array} string
// @Router /currencies/{id} [get]
func GetCurrency(w http.ResponseWriter, r *http.Request) {
	currencyID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(currencyID)
	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &models.Currency{
		ID: uint(id),
	}

	err = models.DB.Model(&models.Currency{}).First(&req).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

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
// @Failure 400 {array} string
// @Router /currencies [post]
func CreateCurrency(w http.ResponseWriter, r *http.Request) {

	req := &models.Currency{}

	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = models.DB.Model(&models.Currency{}).Create(&req).Error

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
// @Failure 400 {array} string
// @Router /currencies/{id} [put]
func UpdateCurrency(w http.ResponseWriter, r *http.Request) {
	currencyID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(currencyID)
	if err != nil {
		validation.InvalidID(w, r)
		return
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
// @Failure 400 {array} string
// @Router /currencies/{id} [delete]
func DeleteCurrency(w http.ResponseWriter, r *http.Request) {
	currencyID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(currencyID)
	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	currency := &models.Currency{
		ID: uint(id),
	}

	// check record exists or not
	err = models.DB.First(&currency).Error
	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	models.DB.Delete(&currency)

	json.NewEncoder(w).Encode(currency)
}
