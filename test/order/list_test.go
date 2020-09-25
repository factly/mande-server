package order

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/dataset"
	"github.com/factly/data-portal-server/test/payment"
	"github.com/factly/data-portal-server/test/product"
	"github.com/factly/data-portal-server/test/tag"
	"github.com/gavv/httpexpect"
)

func TestListOrder(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	adminExpect := httpexpect.New(t, server.URL)

	CommonListTests(t, mock, adminExpect)

	server.Close()

	router = action.RegisterUserRoutes()
	server = httptest.NewServer(router)
	userExpect := httpexpect.New(t, server.URL)

	CommonListTests(t, mock, userExpect)

	server.Close()
}

func CommonListTests(t *testing.T, mock sqlmock.Sqlmock, e *httpexpect.Expect) {
	t.Run("get empty order list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(OrderCols))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_product"`)).
			WillReturnRows(sqlmock.NewRows(append(OrderCols, []string{"product_id", "order_id"}...)))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_dataset"`)).
			WillReturnRows(sqlmock.NewRows(append(OrderCols, []string{"product_id", "dataset_id"}...)))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tag"`)).
			WillReturnRows(sqlmock.NewRows(append(OrderCols, []string{"product_id", "tag_id"}...)))

		e.GET(basePath).
			WithHeader("X-User", "1").
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
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(OrderCols).
				AddRow(1, time.Now(), time.Now(), nil, orderlist[0]["user_id"], orderlist[0]["status"], orderlist[0]["payment_id"], orderlist[0]["razorpay_order_id"]).
				AddRow(2, time.Now(), time.Now(), nil, orderlist[1]["user_id"], orderlist[1]["status"], orderlist[1]["payment_id"], orderlist[1]["razorpay_order_id"]))

		payment.PaymentSelectMock(mock)

		currency.CurrencySelectMock(mock)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_product" INNER JOIN "dp_order_item"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows(append(product.ProductCols, []string{"product_id", "order_id"}...)).
				AddRow(1, time.Now(), time.Now(), nil, product.Product["title"], product.Product["slug"], product.Product["price"], product.Product["status"], product.Product["currency_id"], product.Product["featured_medium_id"], 1, 1))

		dataset.DatasetSelectMock(mock)

		tag.TagSelectMock(mock)

		e.GET(basePath).
			WithHeader("X-User", "1").
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
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(OrderCols).
				AddRow(2, time.Now(), time.Now(), nil, orderlist[1]["user_id"], orderlist[1]["status"], orderlist[1]["payment_id"], orderlist[1]["razorpay_order_id"]))

		payment.PaymentSelectMock(mock)

		currency.CurrencySelectMock(mock)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_product" INNER JOIN "dp_order_item"`)).
			WithArgs(2).
			WillReturnRows(sqlmock.NewRows(append(product.ProductCols, []string{"product_id", "order_id"}...)).
				AddRow(1, time.Now(), time.Now(), nil, product.Product["title"], product.Product["slug"], product.Product["price"], product.Product["status"], product.Product["currency_id"], product.Product["featured_medium_id"], 1, 2))

		dataset.DatasetSelectMock(mock)

		tag.TagSelectMock(mock)

		e.GET(basePath).
			WithHeader("X-User", "1").
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
