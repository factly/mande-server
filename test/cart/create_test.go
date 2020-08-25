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
	"github.com/factly/data-portal-server/test/product"
	"github.com/factly/data-portal-server/test/tag"
	"github.com/gavv/httpexpect"
)

func TestCreateCart(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("create a cart", func(t *testing.T) {
		product.ProductSelectMock(mock)

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_cart"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, Cart["status"], Cart["user_id"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectExec(`INSERT INTO "dp_cart_item"`).
			WithArgs(1, 1, 1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		CartSelectMock(mock)

		productsAssociationSelectMock(mock, 1)

		currency.CurrencySelectMock(mock)

		medium.MediumSelectMock(mock)

		tag.TagSelectMock(mock)

		dataset.DatasetSelectMock(mock)

		e.POST(basePath).
			WithJSON(Cart).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(CartReceive)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable cart body", func(t *testing.T) {
		e.POST(basePath).
			WithJSON(invalidCart).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("empty cart body", func(t *testing.T) {
		e.POST(basePath).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("user does not exist", func(t *testing.T) {
		product.ProductSelectMock(mock)

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_cart"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, Cart["status"], Cart["user_id"]).
			WillReturnError(errCartProductFK)
		mock.ExpectRollback()

		e.POST(basePath).
			WithJSON(Cart).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})
}
