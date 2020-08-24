package catalog

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

func TestCreateCatalog(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("create a catalog", func(t *testing.T) {
		product.ProductSelectMock(mock)

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_catalog"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, Catalog["title"], Catalog["description"], Catalog["featured_medium_id"], test.AnyTime{}).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectExec(`INSERT INTO "dp_catalog_product"`).
			WithArgs(1, 1, 1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		CatalogSelectMock(mock)

		medium.MediumSelectMock(mock)

		productsAssociationSelectMock(mock, 1)

		currency.CurrencySelectMock(mock)

		medium.MediumSelectMock(mock)

		tag.TagSelectMock(mock)

		dataset.DatasetSelectMock(mock)

		result := e.POST(basePath).
			WithJSON(Catalog).
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
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("empty catalog body", func(t *testing.T) {
		e.POST(basePath).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("featured medium does not exist", func(t *testing.T) {
		product.ProductSelectMock(mock)

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_catalog"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, Catalog["title"], Catalog["description"], Catalog["featured_medium_id"], test.AnyTime{}).
			WillReturnError(errCatalogProductFK)
		mock.ExpectRollback()

		e.POST(basePath).
			WithJSON(Catalog).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})
}
