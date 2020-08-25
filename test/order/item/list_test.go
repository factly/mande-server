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

func TestListOrderItem(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("get empty order item list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(OrderItemCols))

		e.GET(basePath).
			WithPath("order_id", "1").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 0})

		test.ExpectationsMet(t, mock)
	})

	t.Run("get order item list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(orderitemlist)))

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(OrderItemCols).
				AddRow(1, time.Now(), time.Now(), nil, orderitemlist[0]["extra_info"], orderitemlist[0]["product_id"], 1).
				AddRow(2, time.Now(), time.Now(), nil, orderitemlist[1]["extra_info"], orderitemlist[1]["product_id"], 1))

		product.ProductSelectMock(mock)

		currency.CurrencySelectMock(mock)

		e.GET(basePath).
			WithPath("order_id", "1").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(orderitemlist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(orderitemlist[0])

		test.ExpectationsMet(t, mock)
	})

	t.Run("get order item list with paiganation", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(orderitemlist)))

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(OrderItemCols).
				AddRow(2, time.Now(), time.Now(), nil, orderitemlist[1]["extra_info"], orderitemlist[1]["product_id"], 1))

		product.ProductSelectMock(mock)

		currency.CurrencySelectMock(mock)

		e.GET(basePath).
			WithPath("order_id", "1").
			WithQueryObject(map[string]interface{}{
				"limit": "1",
				"page":  "2",
			}).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(orderitemlist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(orderitemlist[1])

		test.ExpectationsMet(t, mock)
	})
}
