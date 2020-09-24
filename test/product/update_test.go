package product

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
	"github.com/factly/data-portal-server/test/medium"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestUpdateProduct(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	test.MeiliGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

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

	t.Run("undecodable product body", func(t *testing.T) {
		e.PUT(path).
			WithPath("product_id", "1").
			WithJSON(undecodableProduct).
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

	t.Run("deleting old tags fails", func(t *testing.T) {
		preUpdateMock(mock)

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_product_tag"`)).
			WillReturnError(errors.New("cannot delete tags"))

		mock.ExpectRollback()

		e.PUT(path).
			WithPath("product_id", "1").
			WithJSON(Product).
			Expect().
			Status(http.StatusInternalServerError)
		test.ExpectationsMet(t, mock)
	})

	t.Run("deleting old datasets fails", func(t *testing.T) {
		preUpdateMock(mock)

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_product_tag"`)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_product_dataset"`)).
			WillReturnError(errors.New("cannot delete datasets"))

		mock.ExpectRollback()

		e.PUT(path).
			WithPath("product_id", "1").
			WithJSON(Product).
			Expect().
			Status(http.StatusInternalServerError)
		test.ExpectationsMet(t, mock)
	})

	t.Run("update product with new featured_medium_id = 0", func(t *testing.T) {
		updateMockWithoutMedium(mock)

		ProductSelectMock(mock)

		currency.CurrencySelectMock(mock)

		medium.MediumSelectMock(mock)

		tagsAssociationSelectMock(mock, 1)

		datasetsAssociationSelectMock(mock, 1)

		mock.ExpectCommit()

		Product["featured_medium_id"] = 0
		result := e.PUT(path).
			WithPath("product_id", "1").
			WithJSON(Product).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(ProductReceive)

		validateAssociations(result)

		Product["featured_medium_id"] = 1

		test.ExpectationsMet(t, mock)
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

	t.Run("update product when meili is down", func(t *testing.T) {
		gock.Off()
		updateMock(mock, nil)

		ProductSelectMock(mock)

		currency.CurrencySelectMock(mock)

		medium.MediumSelectMock(mock)

		tagsAssociationSelectMock(mock, 1)

		datasetsAssociationSelectMock(mock, 1)

		mock.ExpectRollback()

		e.PUT(path).
			WithPath("product_id", "1").
			WithJSON(Product).
			Expect().
			Status(http.StatusInternalServerError)
		test.ExpectationsMet(t, mock)
	})

}
