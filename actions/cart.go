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

// Cart request body
type cart struct {
	Status string `json:"status"`
	UserID uint   `json:"user_id"`
}

// GetCart - Get cart by id
// @Summary Show a cart by id
// @Description Get cart by ID
// @Tags Cart
// @ID get-cart-by-id
// @Produce  json
// @Param id path string true "Cart ID"
// @Success 200 {object} models.Cart
// @Failure 400 {array} string
// @Router /carts/{id} [get]
func GetCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cartID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(cartID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &models.Cart{
		ID: uint(id),
	}

	err = models.DB.Model(&models.Cart{}).First(&req).Error
	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(req)
}

// CreateCart - create cart
// @Summary Create cart
// @Description create cart
// @Tags Cart
// @ID add-cart
// @Consume json
// @Produce  json
// @Param Cart body cart true "Cart object"
// @Success 200 {object} models.Cart
// @Failure 400 {array} string
// @Router /carts [post]
func CreateCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	req := &models.Cart{}

	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = models.DB.Model(&models.Cart{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(req)
}

// UpdateCart - Update cart by id
// @Summary Update a cart by id
// @Description Update cart by ID
// @Tags Cart
// @ID update-cart-by-id
// @Produce json
// @Consume json
// @Param id path string true "Cart ID"
// @Param Cart body cart false "Cart"
// @Success 200 {object} models.Cart
// @Failure 400 {array} string
// @Router /carts/{id} [put]
func UpdateCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cartID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(cartID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &models.Cart{}
	cart := &models.Cart{
		ID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&req)

	models.DB.Model(&cart).Updates(models.Cart{
		Status: req.Status,
		UserID: req.UserID,
	})
	models.DB.First(&cart)

	json.NewEncoder(w).Encode(cart)
}

// DeleteCart - Delete cart by id
// @Summary Delete a cart
// @Description Delete cart by ID
// @Tags Cart
// @ID delete-cart-by-id
// @Consume  json
// @Param id path string true "Cart ID"
// @Success 200 {object} models.Cart
// @Failure 400 {array} string
// @Router /carts/{id} [delete]
func DeleteCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cartID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(cartID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	cart := &models.Cart{
		ID: uint(id),
	}

	// check record exists or not
	err = models.DB.First(&cart).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	models.DB.Delete(&cart)

	json.NewEncoder(w).Encode(cart)
}
