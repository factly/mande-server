package payment

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/mande-server/action"
	"github.com/factly/mande-server/test"
	"github.com/factly/mande-server/test/currency"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestListPayment(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	adminExpect := httpexpect.New(t, server.URL)

	test.KavachGock()
	test.KetoGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	CommonListTests(t, mock, adminExpect)

	server.Close()

	router = action.RegisterUserRoutes()
	server = httptest.NewServer(router)
	userExpect := httpexpect.New(t, server.URL)

	gock.New(server.URL).EnableNetworking().Persist()

	CommonListTests(t, mock, userExpect)

	server.Close()
}

func CommonListTests(t *testing.T, mock sqlmock.Sqlmock, e *httpexpect.Expect) {
	t.Run("get empty payments list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(PaymentCols))

		e.GET(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 0})

		test.ExpectationsMet(t, mock)
	})

	t.Run("get payments list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(paymentlist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(PaymentCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, paymentlist[0]["amount"], paymentlist[0]["gateway"], paymentlist[0]["currency_id"], paymentlist[0]["status"], paymentlist[0]["razorpay_payment_id"], paymentlist[0]["razorpay_signature"]).
				AddRow(2, time.Now(), time.Now(), nil, 1, 1, paymentlist[1]["amount"], paymentlist[1]["gateway"], paymentlist[1]["currency_id"], paymentlist[1]["status"], paymentlist[1]["razorpay_payment_id"], paymentlist[1]["razorpay_signature"]))

		currency.CurrencySelectMock(mock)

		delete(paymentlist[0], "razorpay_order_id")
		e.GET(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(paymentlist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(paymentlist[0])

		test.ExpectationsMet(t, mock)
	})

	t.Run("get payments list with paiganation", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(paymentlist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(PaymentCols).
				AddRow(2, time.Now(), time.Now(), nil, 1, 1, paymentlist[1]["amount"], paymentlist[1]["gateway"], paymentlist[1]["currency_id"], paymentlist[1]["status"], paymentlist[1]["razorpay_payment_id"], paymentlist[1]["razorpay_signature"]))

		currency.CurrencySelectMock(mock)

		delete(paymentlist[1], "razorpay_order_id")
		e.GET(basePath).
			WithQueryObject(map[string]interface{}{
				"limit": "1",
				"page":  "2",
			}).
			WithHeaders(headers).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(paymentlist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(paymentlist[1])

		test.ExpectationsMet(t, mock)
	})
}
