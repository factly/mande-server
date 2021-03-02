package order

import (
	"net/http"
	"strconv"

	"github.com/factly/data-portal-server/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/middlewarex"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// userDetails - Get orders by id
// @Summary Show a orders by id
// @Description Get orders by ID
// @Tags Order
// @ID get-orders-by-id
// @Produce  json
// @Param order_id path string true "Order ID"
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Success 200 {object} model.Order
// @Failure 400 {array} string
// @Router /orders/{order_id} [get]
func userDetails(w http.ResponseWriter, r *http.Request) {

	uID, err := middlewarex.GetUser(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	orderID := chi.URLParam(r, "order_id")
	id, err := strconv.Atoi(orderID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &model.Order{}
	result.ID = uint(id)

	err = model.DB.Model(&model.Order{}).Where(&model.Order{
		UserID: uint(uID),
	}).Preload("Payment").Preload("Payment.Currency").Preload("Products").Preload("Products.Datasets").Preload("Products.Tags").First(&result).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	renderx.JSON(w, http.StatusOK, result)
}

// adminDetails - Get orders by id
// @Summary Show a orders by id
// @Description Get orders by ID
// @Tags Order
// @ID get-orders-by-id
// @Produce  json
// @Param order_id path string true "Order ID"
// @Param X-User header string true "User ID"
// @Param X-Organisation header string true "Organisation ID"
// @Success 200 {object} model.Order
// @Failure 400 {array} string
// @Router /orders/{order_id} [get]
func adminDetails(w http.ResponseWriter, r *http.Request) {

	orderID := chi.URLParam(r, "order_id")
	id, err := strconv.Atoi(orderID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &model.Order{}
	result.ID = uint(id)

	err = model.DB.Model(&model.Order{}).Preload("Payment").Preload("Payment.Currency").Preload("Products").Preload("Products.Datasets").Preload("Products.Tags").First(&result).Error

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	renderx.JSON(w, http.StatusOK, result)
}
