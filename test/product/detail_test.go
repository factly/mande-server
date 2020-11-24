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
	"gopkg.in/h2non/gock.v1"
)

func TestDetailProduct(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	adminExpect := httpexpect.New(t, server.URL)

	test.KavachGock()
	test.KetoGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	// ADMIN tests
	adminDetailTests(t, mock, adminExpect)

	server.Close()

	router = action.RegisterUserRoutes()
	server = httptest.NewServer(router)
	userExpect := httpexpect.New(t, server.URL)

	gock.New(server.URL).EnableNetworking().Persist()

	userDetailTests(t, mock, userExpect)

	server.Close()

}

func adminDetailTests(t *testing.T, mock sqlmock.Sqlmock, e *httpexpect.Expect) {
	t.Run("get product by id", func(t *testing.T) {
		ProductSelectMock(mock)

		currency.CurrencySelectMock(mock)

		datasetsAssociationSelectMock(mock, 1)
		medium.MediumSelectMock(mock)

		tagsAssociationSelectMock(mock, 1)

		result := e.GET(path).
			WithPath("product_id", "1").
			WithHeaders(headers).
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

		e.GET(path).
			WithPath("product_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid product id", func(t *testing.T) {
		e.GET(path).
			WithPath("product_id", "abc").
			WithHeaders(headers).
			Expect().
			Status(http.StatusNotFound)
	})
}

func userDetailTests(t *testing.T, mock sqlmock.Sqlmock, e *httpexpect.Expect) {
	t.Run("get product by id", func(t *testing.T) {
		ProductSelectMock(mock)

		currency.CurrencySelectMock(mock)
		datasetsAssociationSelectMock(mock, 1)
		medium.MediumSelectMock(mock)
		tagsAssociationSelectMock(mock, 1)

		catalogsAssociationSelectMock(mock, 1)

		plansCatalogAssociationSelectMock(mock)

		membershipAssociationSelectMock(mock)

		planSelectMock(mock)

		result := e.GET(path).
			WithPath("product_id", "1").
			WithHeaders(headers).
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

		e.GET(path).
			WithPath("product_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid product id", func(t *testing.T) {
		e.GET(path).
			WithPath("product_id", "abc").
			WithHeaders(headers).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("invalid user id header", func(t *testing.T) {
		e.GET(path).
			WithPath("product_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusNotFound)
	})
}
