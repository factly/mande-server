package order

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/mande-server/action"
	"github.com/factly/mande-server/test"
	"github.com/factly/mande-server/test/currency"
	"github.com/factly/mande-server/test/payment"
	"github.com/factly/mande-server/test/product"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestListOrder(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	adminExpect := httpexpect.New(t, server.URL)

	test.KetoGock()
	test.KavachGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	// ADMIN tests
	CommonListTests(t, mock, adminExpect)

	t.Run("get order list with user query parameter", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(orderlist)))

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(OrderCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, orderlist[0]["user_id"], orderlist[0]["status"], orderlist[0]["payment_id"], orderlist[0]["razorpay_order_id"]).
				AddRow(2, time.Now(), time.Now(), nil, 1, 1, orderlist[1]["user_id"], orderlist[1]["status"], orderlist[1]["payment_id"], orderlist[1]["razorpay_order_id"]))

		payment.PaymentSelectMock(mock)

		currency.CurrencySelectMock(mock)

		associatedProductsSelectMock(mock)

		adminExpect.GET(basePath).
			WithHeaders(headers).
			WithQuery("user", "1").
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

	t.Run("invalid user query parameter", func(t *testing.T) {
		adminExpect.GET(basePath).
			WithHeaders(headers).
			WithQuery("user", "abc").
			Expect().
			Status(http.StatusBadRequest)
	})

	server.Close()

	router = action.RegisterUserRoutes()
	server = httptest.NewServer(router)
	userExpect := httpexpect.New(t, server.URL)

	gock.New(server.URL).EnableNetworking().Persist()

	// USER tests
	CommonListTests(t, mock, userExpect)

	server.Close()
}

func CommonListTests(t *testing.T, mock sqlmock.Sqlmock, e *httpexpect.Expect) {
	t.Run("get empty order list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(OrderCols))

		e.GET(basePath).
			WithHeaders(headers).
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
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, orderlist[0]["user_id"], orderlist[0]["status"], orderlist[0]["payment_id"], orderlist[0]["razorpay_order_id"]).
				AddRow(2, time.Now(), time.Now(), nil, 1, 1, orderlist[1]["user_id"], orderlist[1]["status"], orderlist[1]["payment_id"], orderlist[1]["razorpay_order_id"]))

		payment.PaymentSelectMock(mock)

		currency.CurrencySelectMock(mock)

		associatedProductsSelectMock(mock)

		e.GET(basePath).
			WithHeaders(headers).
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
				AddRow(2, time.Now(), time.Now(), nil, 1, 1, orderlist[1]["user_id"], orderlist[1]["status"], orderlist[1]["payment_id"], orderlist[1]["razorpay_order_id"]))

		payment.PaymentSelectMock(mock)

		currency.CurrencySelectMock(mock)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_order_item"`)).
			WithArgs(2).
			WillReturnRows(sqlmock.NewRows([]string{"order_id", "product_id"}).
				AddRow(1, 1))
		product.ProductSelectMock(mock)

		e.GET(basePath).
			WithHeaders(headers).
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
