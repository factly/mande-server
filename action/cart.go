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

// Cart request body
type cart struct {
	Status string `json:"status"`
	UserID uint   `json:"user_id"`
}

// GetCarts - Get all carts
// @Summary Show all carts
// @Description Get all carts
// @Tags Cart
// @ID get-all-carts
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} model.Cart
// @Router /carts [get]
func GetCarts(w http.ResponseWriter, r *http.Request) {

	var carts []model.Cart
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

	model.DB.Offset(offset).Limit(limit).Model(&model.Cart{}).Find(&carts)

	json.NewEncoder(w).Encode(carts)
}

// GetCart - Get cart by id
// @Summary Show a cart by id
// @Description Get cart by ID
// @Tags Cart
// @ID get-cart-by-id
// @Produce  json
// @Param id path string true "Cart ID"
// @Success 200 {object} model.Cart
// @Failure 400 {array} string
// @Router /carts/{id} [get]
func GetCart(w http.ResponseWriter, r *http.Request) {

	cartID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(cartID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &model.Cart{
		ID: uint(id),
	}

	err = model.DB.Model(&model.Cart{}).First(&req).Error
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
// @Success 200 {object} model.Cart
// @Failure 400 {array} string
// @Router /carts [post]
func CreateCart(w http.ResponseWriter, r *http.Request) {

	req := &model.Cart{}

	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = model.DB.Model(&model.Cart{}).Create(&req).Error

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
// @Success 200 {object} model.Cart
// @Failure 400 {array} string
// @Router /carts/{id} [put]
func UpdateCart(w http.ResponseWriter, r *http.Request) {

	cartID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(cartID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &model.Cart{}
	cart := &model.Cart{
		ID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&req)

	model.DB.Model(&cart).Updates(model.Cart{
		Status: req.Status,
		UserID: req.UserID,
	})
	model.DB.First(&cart)

	json.NewEncoder(w).Encode(cart)
}

// DeleteCart - Delete cart by id
// @Summary Delete a cart
// @Description Delete cart by ID
// @Tags Cart
// @ID delete-cart-by-id
// @Consume  json
// @Param id path string true "Cart ID"
// @Success 200 {object} model.Cart
// @Failure 400 {array} string
// @Router /carts/{id} [delete]
func DeleteCart(w http.ResponseWriter, r *http.Request) {

	cartID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(cartID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	cart := &model.Cart{
		ID: uint(id),
	}

	// check record exists or not
	err = model.DB.First(&cart).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	model.DB.Delete(&cart)

	json.NewEncoder(w).Encode(cart)
}
