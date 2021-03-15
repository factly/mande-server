package catalog

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/mande-server/action"
	"github.com/factly/mande-server/test"
	"github.com/factly/mande-server/test/currency"
	"github.com/factly/mande-server/test/medium"
	"github.com/factly/mande-server/test/product"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestCreateCatalog(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	test.MeiliGock()
	test.KetoGock()
	test.KavachGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	e := httpexpect.New(t, server.URL)

	t.Run("create a catalog", func(t *testing.T) {
		product.ProductSelectMock(mock)

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_catalog"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, Catalog["title"], Catalog["description"], test.AnyTime{}, Catalog["featured_medium_id"]).
			WillReturnRows(sqlmock.NewRows([]string{"id", "featured_medium_id"}).AddRow(1, 1))

		mock.ExpectQuery(`INSERT INTO "dp_product"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, product.Product["title"], product.Product["slug"], product.Product["price"], product.Product["status"], product.Product["currency_id"], product.Product["featured_medium_id"], 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectExec(`INSERT INTO "dp_catalog_product"`).
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		CatalogSelectMock(mock)

		medium.MediumSelectMock(mock)
		productsAssociationSelectMock(mock, 1)

		currency.CurrencySelectMock(mock)
		datasetsAssociationSelectMock(mock)
		medium.MediumSelectMock(mock)
		tagsAssociationSelectMock(mock)

		mock.ExpectCommit()

		result := e.POST(basePath).
			WithJSON(Catalog).
			WithHeaders(headers).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(CatalogReceive)

		validateAssociations(result)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable catalog body", func(t *testing.T) {
		e.POST(basePath).
			WithJSON(invalidCatalog).
			WithHeaders(headers).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("empty catalog body", func(t *testing.T) {
		e.POST(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("featured medium does not exist", func(t *testing.T) {
		product.ProductSelectMock(mock)

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_catalog"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, Catalog["title"], Catalog["description"], test.AnyTime{}, Catalog["featured_medium_id"]).
			WillReturnError(errCatalogProductFK)
		mock.ExpectRollback()

		e.POST(basePath).
			WithJSON(Catalog).
			WithHeaders(headers).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("create a catalog when meili is down", func(t *testing.T) {
		gock.Off()
		test.KetoGock()
		test.KavachGock()
		gock.New(server.URL).EnableNetworking().Persist()

		product.ProductSelectMock(mock)

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_catalog"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, Catalog["title"], Catalog["description"], test.AnyTime{}, Catalog["featured_medium_id"]).
			WillReturnRows(sqlmock.NewRows([]string{"id", "featured_medium_id"}).AddRow(1, 1))

		mock.ExpectQuery(`INSERT INTO "dp_product"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, product.Product["title"], product.Product["slug"], product.Product["price"], product.Product["status"], product.Product["currency_id"], product.Product["featured_medium_id"], 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectExec(`INSERT INTO "dp_catalog_product"`).
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		CatalogSelectMock(mock)

		medium.MediumSelectMock(mock)
		productsAssociationSelectMock(mock, 1)

		currency.CurrencySelectMock(mock)
		datasetsAssociationSelectMock(mock)
		medium.MediumSelectMock(mock)
		tagsAssociationSelectMock(mock)

		mock.ExpectRollback()

		e.POST(basePath).
			WithHeaders(headers).
			WithJSON(Catalog).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
