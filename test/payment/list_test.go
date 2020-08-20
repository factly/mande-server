package payment

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/gavv/httpexpect"
)

func TestListPayment(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("get empty payments list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(paymentCols))

		e.GET(basePath).
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
			WillReturnRows(sqlmock.NewRows(paymentCols).
				AddRow(1, time.Now(), time.Now(), nil, paymentlist[0]["amount"], paymentlist[0]["gateway"], paymentlist[0]["currency_id"], paymentlist[0]["status"]).
				AddRow(2, time.Now(), time.Now(), nil, paymentlist[1]["amount"], paymentlist[1]["gateway"], paymentlist[1]["currency_id"], paymentlist[1]["status"]))

		paymentCurrencyMock(mock)

		e.GET(basePath).
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
			WillReturnRows(sqlmock.NewRows(paymentCols).
				AddRow(2, time.Now(), time.Now(), nil, paymentlist[1]["amount"], paymentlist[1]["gateway"], paymentlist[1]["currency_id"], paymentlist[1]["status"]))

		paymentCurrencyMock(mock)

		e.GET(basePath).
			WithQueryObject(map[string]interface{}{
				"limit": "1",
				"page":  "2",
			}).
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
