package product

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/medium"
	"github.com/gavv/httpexpect"
)

func TestUpdateProduct(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("update product", func(t *testing.T) {
		updateMock(mock, nil)

		ProductSelectMock(mock)

		currency.CurrencySelectMock(mock)

		medium.MediumSelectMock(mock)

		tagsAssociationSelectMock(mock, 1)

		datasetsAssociationSelectMock(mock, 1)

		mock.ExpectCommit()

		result := e.PUT(path).
			WithPath("product_id", "1").
			WithJSON(Product).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(ProductReceive)

		validateAssociations(result)

		test.ExpectationsMet(t, mock)
	})

	t.Run("product record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(ProductCols))

		e.PUT(path).
			WithPath("product_id", "1").
			WithJSON(Product).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable product body", func(t *testing.T) {
		e.PUT(path).
			WithPath("product_id", "1").
			WithJSON(invalidProduct).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("invalid product id", func(t *testing.T) {
		e.PUT(path).
			WithPath("product_id", "abc").
			WithJSON(Product).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("new featured medium does not exist", func(t *testing.T) {
		updateMock(mock, errProductMediumFK)

		e.PUT(path).
			WithPath("product_id", "1").
			WithJSON(Product).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("new currency does not exist", func(t *testing.T) {
		updateMock(mock, errProductCurrencyFK)

		e.PUT(path).
			WithPath("product_id", "1").
			WithJSON(Product).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
