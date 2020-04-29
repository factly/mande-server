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

// Order request body
type order struct {
	UserID    uint   `json:"user_id"`
	Status    string `json:"status"`
	PaymentID uint   `json:"payment_id"`
	CartID    uint   `json:"cart_id"`
}

// GetOrders - Get all orders
// @Summary Show all orders
// @Description Get all orders
// @Tags Order
// @ID get-all-orders
// @Produce  json
// @Param limit query string false "limt per page"
// @Param page query string false "page number"
// @Success 200 {array} model.Order
// @Router /orders [get]
func GetOrders(w http.ResponseWriter, r *http.Request) {

	var orders []model.Order
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

	model.DB.Offset(offset).Limit(limit).Preload("Payment").Preload("Payment.Currency").Preload("Cart").Model(&model.Order{}).Find(&orders)

	json.NewEncoder(w).Encode(orders)
}

// GetOrder - Get orders by id
// @Summary Show a orders by id
// @Description Get orders by ID
// @Tags Order
// @ID get-orders-by-id
// @Produce  json
// @Param id path string true "Order ID"
// @Success 200 {object} model.Order
// @Failure 400 {array} string
// @Router /orders/{id} [get]
func GetOrder(w http.ResponseWriter, r *http.Request) {

	orderID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(orderID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &model.Order{
		ID: uint(id),
	}

	err = model.DB.Model(&model.Order{}).Preload("Payment").Preload("Payment.Currency").Preload("Cart").First(&req).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(req)
}

// CreateOrder - create orders
// @Summary Create orders
// @Description create orders
// @Tags Order
// @ID add-orders
// @Consume json
// @Produce  json
// @Param Order body order true "Order object"
// @Success 200 {object} model.Order
// @Failure 400 {array} string
// @Router /orders [post]
func CreateOrder(w http.ResponseWriter, r *http.Request) {

	req := &model.Order{}

	json.NewDecoder(r.Body).Decode(&req)

	validate := validator.New()
	err := validate.StructExcept(req, "Payment", "Cart")
	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	err = model.DB.Model(&model.Order{}).Create(&req).Error

	if err != nil {
		log.Fatal(err)
	}
	model.DB.Model(&model.Order{}).Preload("Payment").Preload("Payment.Currency").Preload("Cart").First(&req)
	json.NewEncoder(w).Encode(req)
}

// UpdateOrder - Update orders by id
// @Summary Update a orders by id
// @Description Update orders by ID
// @Tags Order
// @ID update-orders-by-id
// @Produce json
// @Consume json
// @Param id path string true "Order ID"
// @Param Order body order false "Order"
// @Success 200 {object} model.Order
// @Failure 400 {array} string
// @Router /orders/{id} [put]
func UpdateOrder(w http.ResponseWriter, r *http.Request) {

	orderID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(orderID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	req := &model.Order{}
	orders := &model.Order{
		ID: uint(id),
	}

	json.NewDecoder(r.Body).Decode(&req)

	model.DB.Model(&orders).Updates(model.Order{
		UserID:    req.UserID,
		PaymentID: req.PaymentID,
		Status:    req.Status,
		CartID:    req.CartID,
	})
	model.DB.Preload("Payment").Preload("Payment.Currency").Preload("Cart").First(&orders)

	json.NewEncoder(w).Encode(orders)
}

// DeleteOrder - Delete orders by id
// @Summary Delete a orders
// @Description Delete orders by ID
// @Tags Order
// @ID delete-orders-by-id
// @Consume  json
// @Param id path string true "Order ID"
// @Success 200 {object} model.Order
// @Failure 400 {array} string
// @Router /orders/{id} [delete]
func DeleteOrder(w http.ResponseWriter, r *http.Request) {

	orderID := chi.URLParam(r, "id")
	id, err := strconv.Atoi(orderID)

	if err != nil {
		validation.InvalidID(w, r)
		return
	}

	orders := &model.Order{
		ID: uint(id),
	}

	// check record exists or not
	err = model.DB.First(&orders).Error

	if err != nil {
		validation.RecordNotFound(w, r)
		return
	}
	model.DB.Preload("Payment").Preload("Payment.Currency").Preload("Cart").Delete(&orders)

	json.NewEncoder(w).Encode(orders)
}
