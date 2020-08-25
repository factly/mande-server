package order

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/cart"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/payment"
	"github.com/gavv/httpexpect"
)

func TestListOrder(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("get empty order list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(OrderCols))

		e.GET(basePath).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 0})

		test.ExpectationsMet(t, mock)
	})

	t.Run("get order list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(orderlist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(OrderCols).
				AddRow(1, time.Now(), time.Now(), nil, orderlist[0]["user_id"], orderlist[0]["status"], orderlist[0]["payment_id"], orderlist[0]["cart_id"]).
				AddRow(2, time.Now(), time.Now(), nil, orderlist[1]["user_id"], orderlist[1]["status"], orderlist[1]["payment_id"], orderlist[1]["cart_id"]))

		payment.PaymentSelectMock(mock)

		currency.CurrencySelectMock(mock)

		cart.CartSelectMock(mock)

		e.GET(basePath).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(orderlist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(orderlist[0])

		test.ExpectationsMet(t, mock)
	})

	t.Run("get order list with paiganation", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(orderlist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(OrderCols).
				AddRow(2, time.Now(), time.Now(), nil, orderlist[1]["user_id"], orderlist[1]["status"], orderlist[1]["payment_id"], orderlist[1]["cart_id"]))

		payment.PaymentSelectMock(mock)

		currency.CurrencySelectMock(mock)

		cart.CartSelectMock(mock)

		e.GET(basePath).
			WithQueryObject(map[string]interface{}{
				"limit": "1",
				"page":  "2",
			}).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(orderlist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(orderlist[1])

		test.ExpectationsMet(t, mock)
	})
}
