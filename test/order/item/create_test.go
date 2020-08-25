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

func TestCreateOrderItem(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("create an order item", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_order_item"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, OrderItem["extra_info"], OrderItem["product_id"], 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		OrderItemSelectMock(mock)

		product.ProductSelectMock(mock)

		currency.CurrencySelectMock(mock)

		e.POST(basePath).
			WithPath("order_id", "1").
			WithJSON(OrderItem).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(OrderItem).
			Value("product").
			Object().
			ContainsMap(product.ProductReceive)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable order item body", func(t *testing.T) {
		e.POST(basePath).
			WithPath("order_id", "1").
			WithJSON(invalidOrderItem).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("empty order item body", func(t *testing.T) {
		e.POST(basePath).
			WithPath("order_id", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("product does not exist", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_order_item"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, OrderItem["extra_info"], OrderItem["product_id"], 1).
			WillReturnError(errOrderItemFK)
		mock.ExpectRollback()

		e.POST(basePath).
			WithPath("order_id", "1").
			WithJSON(OrderItem).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})
}
