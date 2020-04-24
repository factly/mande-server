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

// CartItem request body
type cartItem struct {
	IsDeleted bool `json:"is_deleted"`
	CartID    uint `json:"cart_id"`
	ProductID uint `json:"product_id"`
}

// GetCartItems - Get all cartItems
// @Summary Show all cartItems
// @Description Get all cartItems
// @Tags CartItem
// @ID get-all-cartItems
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} models.CartItem
// @Router /cartItems [get]
func GetCartItems(w http.ResponseWriter, r *http.Request) {

	var cartItems []models.CartItem
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

	models.DB.Offset(offset).Limit(limit).Preload("Product").Preload("Product.Status").Preload("Product.ProductType").Preload("Product.Currency").Model(&models.CartItem{}).Find(&cartItems)

	json.NewEncoder(w).Encode(cartItems)
}

// GetCartItem - Get cartItem by id
// @Summary Show a cartItem by id
// @Description Get cartItem by ID
// @Tags CartItem
// @ID get-cartItem-by-id
// @Produce  json
// @Param id path string true "CartItem ID"
// @Success 200 {object} models.CartItem
// @Failure 400 {array} string
// @Router /cartItems/{id} [get]
func GetCartItem(w http.ResponseWriter, r *http.Request) {

	cartItemID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(cartItemID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &models.CartItem{
		ID: uint(id),
	}

	err = models.DB.Model(&models.CartItem{}).First(&req).Error
	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	models.DB.Model(&req).Association("Product").Find(&req.Product)
	models.DB.Model(&req.Product).Association("Status").Find(&req.Product.Status)
	models.DB.Model(&req.Product).Association("ProductType").Find(&req.Product.ProductType)
	models.DB.Model(&req.Product).Association("Currency").Find(&req.Product.Currency)
	json.NewEncoder(w).Encode(req)
}

// CreateCartItem - create cartItem
// @Summary Create cartItem
// @Description create cartItem
// @Tags CartItem
// @ID add-cartItem
// @Consume json
// @Produce  json
// @Param CartItem body cartItem true "CartItem object"
// @Success 200 {object} models.CartItem
// @Failure 400 {array} string
// @Router /cartItems [post]
func CreateCartItem(w http.ResponseWriter, r *http.Request) {

	req := &models.CartItem{}

	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.StructExcept(req, "Product")
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}
	err = models.DB.Model(&models.CartItem{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}
	models.DB.Model(&req).Association("Product").Find(&req.Product)
	models.DB.Model(&req.Product).Association("Status").Find(&req.Product.Status)
	models.DB.Model(&req.Product).Association("ProductType").Find(&req.Product.ProductType)
	models.DB.Model(&req.Product).Association("Currency").Find(&req.Product.Currency)
	json.NewEncoder(w).Encode(req)
}

// UpdateCartItem - Update cartItem by id
// @Summary Update a cartItem by id
// @Description Update cartItem by ID
// @Tags CartItem
// @ID update-cartItem-by-id
// @Produce json
// @Consume json
// @Param id path string true "CartItem ID"
// @Param CartItem body cartItem false "CartItem"
// @Success 200 {object} models.CartItem
// @Failure 400 {array} string
// @Router /cartItems/{id} [put]
func UpdateCartItem(w http.ResponseWriter, r *http.Request) {

	cartItemID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(cartItemID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &models.CartItem{}
	cartItem := &models.CartItem{
		ID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&req)

	models.DB.Model(&cartItem).Updates(models.CartItem{
		IsDeleted: req.IsDeleted,
		CartID:    req.CartID,
		ProductID: req.ProductID,
	})
	models.DB.First(&cartItem)
	models.DB.Model(&cartItem).Association("Product").Find(&cartItem.Product)
	models.DB.Model(&cartItem.Product).Association("Status").Find(&cartItem.Product.Status)
	models.DB.Model(&cartItem.Product).Association("ProductType").Find(&cartItem.Product.ProductType)
	models.DB.Model(&cartItem.Product).Association("Currency").Find(&cartItem.Product.Currency)
	json.NewEncoder(w).Encode(cartItem)
}

// DeleteCartItem - Delete cartItem by id
// @Summary Delete a cartItem
// @Description Delete cartItem by ID
// @Tags CartItem
// @ID delete-cartItem-by-id
// @Consume  json
// @Param id path string true "CartItem ID"
// @Success 200 {object} models.CartItem
// @Failure 400 {array} string
// @Router /cartItems/{id} [delete]
func DeleteCartItem(w http.ResponseWriter, r *http.Request) {

	cartItemID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(cartItemID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	cartItem := &models.CartItem{
		ID: uint(id),
	}

	// check record exists or not
	err = models.DB.First(&cartItem).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	models.DB.Delete(&cartItem)
	models.DB.Model(&cartItem).Association("Product").Find(&cartItem.Product)
	models.DB.Model(&cartItem.Product).Association("Status").Find(&cartItem.Product.Status)
	models.DB.Model(&cartItem.Product).Association("ProductType").Find(&cartItem.Product.ProductType)
	models.DB.Model(&cartItem.Product).Association("Currency").Find(&cartItem.Product.Currency)
	json.NewEncoder(w).Encode(cartItem)
}
