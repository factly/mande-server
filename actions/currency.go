package actions

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-api/models"
	"github.com/factly/data-portal-api/validationerrors"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
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
// @Param id path string false "Currency ID"
// @Success 200 {object} models.Currency
// @Failure 400 {array} string
// @Router /currencies/{id} [get]
func GetCurrency(w http.ResponseWriter, r *http.Request) {
	currencyID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(currencyID)
	if err != nil {
		validationerrors.InvalidID(w, r)
		return
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
// @Failure 400 {array} string
// @Router /currencies [post]
func CreateCurrency(w http.ResponseWriter, r *http.Request) {

	req := &models.Currency{}

	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		msg := err.Error()
		validationerrors.ValidErrors(w, r, msg)
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
		validationerrors.InvalidID(w, r)
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
		validationerrors.InvalidID(w, r)
		return
	}

	currency := &models.Currency{
		ID: uint(id),
	}

	// check record exists or not
	err = models.DB.First(&currency).Error
	if err != nil {
		validationerrors.RecordNotFound(w, r)
		return
	}

	models.DB.Delete(&currency)

	json.NewEncoder(w).Encode(currency)
}
