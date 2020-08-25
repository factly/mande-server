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

func TestDetailCart(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("get cart by id", func(t *testing.T) {
		CartSelectMock(mock)

		productsAssociationSelectMock(mock, 1)

		currency.CurrencySelectMock(mock)

		medium.MediumSelectMock(mock)

		tag.TagSelectMock(mock)

		dataset.DatasetSelectMock(mock)

		e.GET(path).
			WithPath("cart_id", "1").
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

		e.GET(path).
			WithPath("cart_id", "1").
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid cart id", func(t *testing.T) {
		e.GET(path).
			WithPath("cart_id", "abc").
			Expect().
			Status(http.StatusNotFound)
	})
}
