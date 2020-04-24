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

// OrderItem request body
type orderItem struct {
	ExtraInfo string `json:"extra_info"`
	ProductID uint   `json:"product_id"`
	OrderID   uint   `json:"order_id"`
}

// GetOrderItems - Get all order items
// @Summary Show all order items
// @Description Get all order items
// @Tags OrderItem
// @ID get-all-order-items
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} models.OrderItem
// @Router /order-items [get]
func GetOrderItems(w http.ResponseWriter, r *http.Request) {

	var orderItems []models.OrderItem
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

	models.DB.Offset(offset).Limit(limit).Preload("Product").Preload("Product.Status").Preload("Product.ProductType").Preload("Product.Currency").Preload("Order").Model(&models.OrderItem{}).Find(&orderItems)

	json.NewEncoder(w).Encode(orderItems)
}

// GetOrderItem - Get order item by id
// @Summary Show a order item by id
// @Description Get order item by ID
// @Tags OrderItem
// @ID get-order-item-by-id
// @Produce  json
// @Param id path string true "Order item ID"
// @Success 200 {object} models.OrderItem
// @Failure 400 {array} string
// @Router /order-items/{id} [get]
func GetOrderItem(w http.ResponseWriter, r *http.Request) {

	orderItemID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(orderItemID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &models.OrderItem{
		ID: uint(id),
	}

	err = models.DB.Model(&models.OrderItem{}).Preload("Product").Preload("Product.Status").Preload("Product.ProductType").Preload("Product.Currency").Preload("Order").First(&req).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(req)
}

// CreateOrderItem - create order items
// @Summary Create order items
// @Description create order items
// @Tags OrderItem
// @ID add-order-item
// @Consume json
// @Produce  json
// @Param OrderItem body orderItem true "Order item object"
// @Success 200 {object} models.OrderItem
// @Router /order-items [post]
func CreateOrderItem(w http.ResponseWriter, r *http.Request) {

	req := &models.OrderItem{}

	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.StructExcept(req, "Product", "Order")
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = models.DB.Model(&models.OrderItem{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}
	models.DB.Model(&models.OrderItem{}).Preload("Product").Preload("Product.Status").Preload("Product.ProductType").Preload("Product.Currency").Preload("Order").First(&req)
	json.NewEncoder(w).Encode(req)
}

// UpdateOrderItem - Update orderItems by id
// @Summary Update a orderItems by id
// @Description Update orderItems by ID
// @Tags OrderItem
// @ID update-orderItems-by-id
// @Produce json
// @Consume json
// @Param id path string true "OrderItem ID"
// @Param OrderItem body orderItem false "OrderItem"
// @Success 200 {object} models.OrderItem
// @Failure 400 {array} string
// @Router /orderItems/{id} [put]
func UpdateOrderItem(w http.ResponseWriter, r *http.Request) {

	orderItemID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(orderItemID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &models.OrderItem{}
	orderItem := &models.OrderItem{
		ID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&req)

	models.DB.Model(&orderItem).Updates(models.OrderItem{
		ExtraInfo: req.ExtraInfo,
		ProductID: req.ProductID,
		OrderID:   req.OrderID,
	})
	models.DB.Preload("Product").Preload("Product.Status").Preload("Product.ProductType").Preload("Product.Currency").Preload("Order").First(&orderItem)

	json.NewEncoder(w).Encode(orderItem)
}

// DeleteOrderItem - Delete order item by id
// @Summary Delete a order item
// @Description Delete order item by ID
// @Tags OrderItem
// @ID delete-order-items-by-id
// @Consume  json
// @Param id path string true "OrderItem ID"
// @Success 200 {object} models.OrderItem
// @Failure 400 {array} string
// @Router /order-items/{id} [delete]
func DeleteOrderItem(w http.ResponseWriter, r *http.Request) {

	orderItemID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(orderItemID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	orderItems := &models.OrderItem{
		ID: uint(id),
	}

	// check record exists or not
	err = models.DB.First(&orderItems).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	models.DB.Preload("Product").Preload("Product.Status").Preload("Product.ProductType").Preload("Product.Currency").Preload("Order").Delete(&orderItems)

	json.NewEncoder(w).Encode(orderItems)
}
