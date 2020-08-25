package cart

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/dataset"
	"github.com/factly/data-portal-server/test/medium"
	"github.com/factly/data-portal-server/test/tag"
	"github.com/gavv/httpexpect"
)

func TestUpdateCart(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("update cart", func(t *testing.T) {
		updateMock(mock, nil)

		CartSelectMock(mock)

		productsAssociationSelectMock(mock, 1)

		currency.CurrencySelectMock(mock)

		medium.MediumSelectMock(mock)

		tag.TagSelectMock(mock)

		dataset.DatasetSelectMock(mock)

		mock.ExpectCommit()

		e.PUT(path).
			WithPath("cart_id", "1").
			WithJSON(Cart).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(CartReceive)

		test.ExpectationsMet(t, mock)
	})

	t.Run("cart record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(CartCols))

		e.PUT(path).
			WithPath("cart_id", "1").
			WithJSON(Cart).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable cart body", func(t *testing.T) {
		e.PUT(path).
			WithPath("cart_id", "1").
			WithJSON(invalidCart).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("invalid cart id", func(t *testing.T) {
		e.PUT(path).
			WithPath("cart_id", "abc").
			WithJSON(Cart).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("new user does not exist", func(t *testing.T) {
		updateMock(mock, errCartProductFK)

		e.PUT(path).
			WithPath("cart_id", "1").
			WithJSON(Cart).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})
}
