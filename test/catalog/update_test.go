package catalog

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
	"github.com/factly/data-portal-server/test/product"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestUpdateCatalog(t *testing.T) {
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

	t.Run("update catalog", func(t *testing.T) {
		updateMock(mock, nil)

		CatalogSelectMock(mock)

		medium.MediumSelectMock(mock)
		productsAssociationSelectMock(mock, 1)

		currency.CurrencySelectMock(mock)
		datasetsAssociationSelectMock(mock)
		medium.MediumSelectMock(mock)
		tagsAssociationSelectMock(mock)

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

	t.Run("undecodable catalog body", func(t *testing.T) {
		e.PUT(path).
			WithPath("catalog_id", "1").
			WithJSON(undecodableCatalog).
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

		test.ExpectationsMet(t, mock)
	})

	t.Run("replacing old products fails", func(t *testing.T) {
		preUpdateMock(mock)

		mock.ExpectQuery(`INSERT INTO "dp_product"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, product.Product["title"], product.Product["slug"], product.Product["price"], product.Product["status"], product.Product["currency_id"], product.Product["featured_medium_id"], 1).
			WillReturnError(errors.New("cannot replace products"))

		mock.ExpectExec(`INSERT INTO "dp_catalog_product"`).
			WithArgs(1, 1).
			WillReturnError(errors.New("cannot replace products"))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_catalog_product"`)).
			WillReturnError(errors.New("cannot replace products"))

		mock.ExpectRollback()

		e.PUT(path).
			WithPath("catalog_id", "1").
			WithJSON(Catalog).
			Expect().
			Status(http.StatusInternalServerError)
		test.ExpectationsMet(t, mock)
	})

	t.Run("update catalog with featured_medium_id = 0", func(t *testing.T) {
		updateWithoutFeaturedMedium(mock)

		CatalogSelectMock(mock)

		medium.MediumSelectMock(mock)
		productsAssociationSelectMock(mock, 1)

		currency.CurrencySelectMock(mock)
		datasetsAssociationSelectMock(mock)
		medium.MediumSelectMock(mock)
		tagsAssociationSelectMock(mock)
		mock.ExpectCommit()

		Catalog["featured_medium_id"] = 0
		result := e.PUT(path).
			WithPath("catalog_id", "1").
			WithJSON(Catalog).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(CatalogReceive)

		validateAssociations(result)
		Catalog["featured_medium_id"] = 1
		test.ExpectationsMet(t, mock)
	})

	t.Run("update catalog when meili is down", func(t *testing.T) {
		gock.Off()
		updateMock(mock, nil)

		CatalogSelectMock(mock)

		medium.MediumSelectMock(mock)
		productsAssociationSelectMock(mock, 1)

		currency.CurrencySelectMock(mock)
		datasetsAssociationSelectMock(mock)
		medium.MediumSelectMock(mock)
		tagsAssociationSelectMock(mock)

		mock.ExpectRollback()

		e.PUT(path).
			WithPath("catalog_id", "1").
			WithJSON(Catalog).
			Expect().
			Status(http.StatusInternalServerError)
		test.ExpectationsMet(t, mock)
	})

}
