package item

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
	"github.com/go-chi/chi"
)

// create - create order items
// @Summary Create order items
// @Description create order items
// @Tags OrderItem
// @ID add-order-item
// @Consume json
// @Produce  json
// @Param order_id path string true "Order ID"
// @Param OrderItem body orderItem true "Order item object"
// @Success 201 {object} model.OrderItem
// @Router /orders/{order_id}/items [post]
func create(w http.ResponseWriter, r *http.Request) {

	orderID := chi.URLParam(r, "order_id")
	id, _ := strconv.Atoi(orderID)

	orderItem := &orderItem{}
	result := &model.OrderItem{}

	err := json.NewDecoder(r.Body).Decode(&orderItem)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	validationError := validationx.Check(orderItem)
	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	result.OrderID = uint(id)
	result.ExtraInfo = orderItem.ExtraInfo
	result.ProductID = orderItem.ProductID

	err = model.DB.Model(&model.OrderItem{}).Create(&result).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}
	model.DB.Model(&model.OrderItem{}).Preload("Product").Preload("Product.Currency").First(&result)

	renderx.JSON(w, http.StatusCreated, result)
}
