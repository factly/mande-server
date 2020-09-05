package orderitem

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/product"
	"github.com/gavv/httpexpect"
)

func TestUpdateOrderItem(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("update order item", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(OrderItemCols).
				AddRow(1, time.Now(), time.Now(), nil, "extra_info", 2, 2))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE \"dp_order_item\" SET (.+)  WHERE (.+) \"dp_order_item\".\"id\" = `).
			WithArgs(OrderItem["extra_info"], OrderItem["product_id"], test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		OrderItemSelectMock(mock)

		product.ProductSelectMock(mock)

		currency.CurrencySelectMock(mock)

		e.PUT(path).
			WithPathObject(map[string]interface{}{
				"order_id": "1",
				"item_id":  "1",
			}).
			WithJSON(OrderItem).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(OrderItem).
			Value("product").
			Object().
			ContainsMap(product.ProductReceive)

		test.ExpectationsMet(t, mock)
	})

	t.Run("order item record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(OrderItemCols))

		e.PUT(path).
			WithPathObject(map[string]interface{}{
				"order_id": "1",
				"item_id":  "1",
			}).
			WithJSON(OrderItem).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid order item id", func(t *testing.T) {
		e.PUT(path).
			WithPathObject(map[string]interface{}{
				"order_id": "1",
				"item_id":  "abc",
			}).
			WithJSON(OrderItem).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("invalid order id", func(t *testing.T) {
		e.PUT(path).
			WithPathObject(map[string]interface{}{
				"order_id": "abc",
				"item_id":  "1",
			}).
			WithJSON(OrderItem).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("invalid order item body", func(t *testing.T) {
		e.PUT(path).
			WithPathObject(map[string]interface{}{
				"order_id": "1",
				"item_id":  "1",
			}).
			WithJSON(invalidOrderItem).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("undecodable order item body", func(t *testing.T) {
		e.PUT(path).
			WithPathObject(map[string]interface{}{
				"order_id": "1",
				"item_id":  "1",
			}).
			WithJSON(undecodableOrderItem).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("new product not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(OrderItemCols).
				AddRow(1, time.Now(), time.Now(), nil, "extra_info", 2, 2))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE \"dp_order_item\" SET (.+)  WHERE (.+) \"dp_order_item\".\"id\" = `).
			WithArgs(OrderItem["extra_info"], OrderItem["product_id"], test.AnyTime{}, 1).
			WillReturnError(errOrderItemFK)
		mock.ExpectRollback()

		e.PUT(path).
			WithPathObject(map[string]interface{}{
				"order_id": "1",
				"item_id":  "1",
			}).
			WithJSON(OrderItem).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})
}
