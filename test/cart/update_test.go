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

	t.Run("update cart item", func(t *testing.T) {
		updateMock(mock, nil)

		CartItemSelectMock(mock)

		product.ProductSelectMock(mock)

		currency.CurrencySelectMock(mock)

		medium.MediumSelectMock(mock)

		tag.TagSelectMock(mock)

		dataset.DatasetSelectMock(mock)

		mock.ExpectCommit()

		e.PUT(path).
			WithHeader("X-User", "1").
			WithPath("cartitem_id", "1").
			WithJSON(CartItem).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(CartItem)

		test.ExpectationsMet(t, mock)
	})

	t.Run("cart item record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(CartItemCols))

		e.PUT(path).
			WithHeader("X-User", "1").
			WithPath("cartitem_id", "1").
			WithJSON(CartItem).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable cart item body", func(t *testing.T) {
		e.PUT(path).
			WithHeader("X-User", "1").
			WithPath("cartitem_id", "1").
			WithJSON(invalidCartItem).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("undecodable cart item body", func(t *testing.T) {
		e.PUT(path).
			WithHeader("X-User", "1").
			WithPath("cartitem_id", "1").
			WithJSON(undecodableCartItem).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("invalid cart item id", func(t *testing.T) {
		e.PUT(path).
			WithHeader("X-User", "1").
			WithPath("cartitem_id", "abc").
			WithJSON(CartItem).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("invalid user header", func(t *testing.T) {
		e.POST(basePath).
			WithHeader("X-User", "abc").
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("new product does not exist", func(t *testing.T) {
		updateMock(mock, errCartItemProductFK)

		e.PUT(path).
			WithHeader("X-User", "1").
			WithPath("cartitem_id", "1").
			WithJSON(CartItem).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("update cart item when meili is down", func(t *testing.T) {
		gock.Off()
		updateMock(mock, nil)

		CartItemSelectMock(mock)

		product.ProductSelectMock(mock)

		currency.CurrencySelectMock(mock)

		medium.MediumSelectMock(mock)

		tag.TagSelectMock(mock)

		dataset.DatasetSelectMock(mock)

		mock.ExpectRollback()

		e.PUT(path).
			WithHeader("X-User", "1").
			WithPath("cartitem_id", "1").
			WithJSON(CartItem).
			Expect().
			Status(http.StatusInternalServerError)
		test.ExpectationsMet(t, mock)
	})

}
