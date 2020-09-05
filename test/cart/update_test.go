package cart

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/dataset"
	"github.com/factly/data-portal-server/test/medium"
	"github.com/factly/data-portal-server/test/tag"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestUpdateCart(t *testing.T) {
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

	t.Run("undecodable cart body", func(t *testing.T) {
		e.PUT(path).
			WithPath("cart_id", "1").
			WithJSON(undecodableCart).
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

	t.Run("deleting old products fails", func(t *testing.T) {
		preUpdateMock(mock)

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_cart_item"`)).
			WillReturnError(errors.New("cannot delete products"))

		mock.ExpectRollback()

		e.PUT(path).
			WithPath("cart_id", "1").
			WithJSON(Cart).
			Expect().
			Status(http.StatusInternalServerError)
		test.ExpectationsMet(t, mock)
	})

	t.Run("update cart when meili is down", func(t *testing.T) {
		gock.Off()
		updateMock(mock, nil)

		CartSelectMock(mock)

		productsAssociationSelectMock(mock, 1)

		currency.CurrencySelectMock(mock)

		medium.MediumSelectMock(mock)

		tag.TagSelectMock(mock)

		dataset.DatasetSelectMock(mock)

		mock.ExpectRollback()

		e.PUT(path).
			WithPath("cart_id", "1").
			WithJSON(Cart).
			Expect().
			Status(http.StatusInternalServerError)
		test.ExpectationsMet(t, mock)
	})

}
