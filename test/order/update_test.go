package order

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/cart"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/payment"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestUpdateOrder(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	test.MeiliGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	e := httpexpect.New(t, server.URL)

	t.Run("update order", func(t *testing.T) {
		updateMock(mock, nil)

		OrderSelectMock(mock)

		payment.PaymentSelectMock(mock)

		currency.CurrencySelectMock(mock)

		cart.CartSelectMock(mock)

		mock.ExpectCommit()

		result := e.PUT(path).
			WithPath("order_id", "1").
			WithJSON(Order).
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

		e.PUT(path).
			WithPath("order_id", "1").
			WithJSON(Order).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable order body", func(t *testing.T) {
		e.PUT(path).
			WithPath("order_id", "1").
			WithJSON(invalidOrder).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("undecodable order body", func(t *testing.T) {
		e.PUT(path).
			WithPath("order_id", "1").
			WithJSON(undecodableOrder).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("invalid order id", func(t *testing.T) {
		e.PUT(path).
			WithPath("order_id", "abc").
			WithJSON(Order).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("new payment does not exist", func(t *testing.T) {
		updateMock(mock, errOrderPaymentFK)
		mock.ExpectRollback()

		e.PUT(path).
			WithPath("order_id", "1").
			WithJSON(Order).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("new user does not exist", func(t *testing.T) {
		updateMock(mock, errOrderUserFK)
		mock.ExpectRollback()

		e.PUT(path).
			WithPath("order_id", "1").
			WithJSON(Order).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("new cart does not exist", func(t *testing.T) {
		updateMock(mock, errOrderCartFK)
		mock.ExpectRollback()

		e.PUT(path).
			WithPath("order_id", "1").
			WithJSON(Order).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("update order when meili is down", func(t *testing.T) {
		gock.Off()
		updateMock(mock, nil)

		OrderSelectMock(mock)

		payment.PaymentSelectMock(mock)

		currency.CurrencySelectMock(mock)

		cart.CartSelectMock(mock)

		mock.ExpectRollback()

		e.PUT(path).
			WithPath("order_id", "1").
			WithJSON(Order).
			Expect().
			Status(http.StatusInternalServerError)
		test.ExpectationsMet(t, mock)
	})

}
