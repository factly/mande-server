package orderitem

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/product"
	"github.com/gavv/httpexpect"
)

func TestDetailOrderItem(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("get order item by id", func(t *testing.T) {
		OrderItemSelectMock(mock)

		product.ProductSelectMock(mock)

		currency.CurrencySelectMock(mock)

		e.GET(path).
			WithPathObject(map[string]interface{}{
				"order_id": "1",
				"item_id":  "1",
			}).
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

	t.Run("order item not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(OrderItemCols))

		e.GET(path).
			WithPathObject(map[string]interface{}{
				"order_id": "1",
				"item_id":  "1",
			}).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid order item id", func(t *testing.T) {
		e.GET(path).
			WithPathObject(map[string]interface{}{
				"order_id": "1",
				"item_id":  "abc",
			}).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("invalid order id", func(t *testing.T) {
		e.GET(path).
			WithPathObject(map[string]interface{}{
				"order_id": "abc",
				"item_id":  "1",
			}).
			Expect().
			Status(http.StatusNotFound)
	})
}
