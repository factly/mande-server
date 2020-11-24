package order

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/payment"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestDetailOrder(t *testing.T) {

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
	CommonDetailTests(t, mock, adminExpect)

	t.Run("get order by id", func(t *testing.T) {
		OrderSelectMock(mock)

		payment.PaymentSelectMock(mock)

		currency.CurrencySelectMock(mock)

		associatedProductsSelectMock(mock)

		result := adminExpect.GET(path).
			WithHeaders(headers).
			WithPath("order_id", "1").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(Order)

		validateAssociations(result)

		test.ExpectationsMet(t, mock)
	})

	t.Run("order record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(OrderCols))

		adminExpect.GET(path).
			WithHeaders(headers).
			WithPath("order_id", "1").
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	server.Close()

	router = action.RegisterUserRoutes()
	server = httptest.NewServer(router)
	userExpect := httpexpect.New(t, server.URL)

	gock.New(server.URL).EnableNetworking().Persist()

	// USER tests
	CommonDetailTests(t, mock, userExpect)

	t.Run("get order by id", func(t *testing.T) {
		OrderSelectMock(mock, 1, 1)

		payment.PaymentSelectMock(mock)

		currency.CurrencySelectMock(mock)

		associatedProductsSelectMock(mock)

		result := userExpect.GET(path).
			WithHeaders(headers).
			WithPath("order_id", "1").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(Order)

		validateAssociations(result)

		test.ExpectationsMet(t, mock)
	})

	t.Run("order record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows(OrderCols))

		userExpect.GET(path).
			WithHeaders(headers).
			WithPath("order_id", "1").
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	server.Close()
}

func CommonDetailTests(t *testing.T, mock sqlmock.Sqlmock, e *httpexpect.Expect) {

	t.Run("invalid order id", func(t *testing.T) {
		e.GET(path).
			WithHeaders(headers).
			WithPath("order_id", "abc").
			Expect().
			Status(http.StatusNotFound)
	})

}
