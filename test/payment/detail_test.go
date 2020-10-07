package payment

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/gavv/httpexpect"
)

func TestDetailPayment(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	adminExpect := httpexpect.New(t, server.URL)

	CommonDetailTests(t, mock, adminExpect)

	server.Close()

	router = action.RegisterUserRoutes()
	server = httptest.NewServer(router)
	userExpect := httpexpect.New(t, server.URL)

	CommonDetailTests(t, mock, userExpect)

	server.Close()
}

func CommonDetailTests(t *testing.T, mock sqlmock.Sqlmock, e *httpexpect.Expect) {
	t.Run("get payment by id", func(t *testing.T) {
		PaymentSelectMock(mock)

		currency.CurrencySelectMock(mock)

		e.GET(path).
			WithPath("payment_id", "1").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(Payment)

		test.ExpectationsMet(t, mock)
	})

	t.Run("payment record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(PaymentCols))

		e.GET(path).
			WithPath("payment_id", "1").
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid payment id", func(t *testing.T) {
		e.GET(path).
			WithPath("payment_id", "abc").
			Expect().
			Status(http.StatusNotFound)
	})
}
