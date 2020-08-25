package order

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/cart"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/payment"
	"github.com/gavv/httpexpect"
)

func TestCreateOrder(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("create a order", func(t *testing.T) {
		insertMock(mock, nil)

		OrderSelectMock(mock)

		payment.PaymentSelectMock(mock)

		currency.CurrencySelectMock(mock)

		cart.CartSelectMock(mock)

		result := e.POST(basePath).
			WithJSON(Order).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(Order)

		validateAssociations(result)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable order body", func(t *testing.T) {
		e.POST(basePath).
			WithJSON(invalidOrder).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("empty order body", func(t *testing.T) {
		e.POST(basePath).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("payment does not exist", func(t *testing.T) {
		insertMock(mock, errOrderPaymentFK)

		e.POST(basePath).
			WithJSON(Order).
			Expect().
			Status(http.StatusInternalServerError)
	})

	t.Run("user does not exist", func(t *testing.T) {
		insertMock(mock, errOrderUserFK)

		e.POST(basePath).
			WithJSON(Order).
			Expect().
			Status(http.StatusInternalServerError)
	})

	t.Run("cart does not exist", func(t *testing.T) {
		insertMock(mock, errOrderCartFK)

		e.POST(basePath).
			WithJSON(Order).
			Expect().
			Status(http.StatusInternalServerError)
	})
}
