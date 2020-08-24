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
	"github.com/factly/data-portal-server/test/tag"
	"github.com/gavv/httpexpect"
)

func TestUpdateCatalog(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("update catalog", func(t *testing.T) {
		updateMock(mock, nil)

		CatalogSelectMock(mock)

		medium.MediumSelectMock(mock)

		productsAssociationSelectMock(mock, 1)

		currency.CurrencySelectMock(mock)

		medium.MediumSelectMock(mock)

		tag.TagSelectMock(mock)

		dataset.DatasetSelectMock(mock)

		mock.ExpectCommit()

		result := e.PUT(path).
			WithPath("catalog_id", "1").
			WithJSON(Catalog).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(CatalogReceive)

		validateAssociations(result)

		test.ExpectationsMet(t, mock)
	})

	t.Run("catalog record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(CatalogCols))

		e.PUT(path).
			WithPath("catalog_id", "1").
			WithJSON(Catalog).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable catalog body", func(t *testing.T) {
		e.PUT(path).
			WithPath("catalog_id", "1").
			WithJSON(invalidCatalog).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("invalid catalog id", func(t *testing.T) {
		e.PUT(path).
			WithPath("catalog_id", "abc").
			WithJSON(Catalog).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("new featured medium does not exist", func(t *testing.T) {
		updateMock(mock, errCatalogProductFK)

		e.PUT(path).
			WithPath("catalog_id", "1").
			WithJSON(Catalog).
			Expect().
			Status(http.StatusInternalServerError)
	})
}
