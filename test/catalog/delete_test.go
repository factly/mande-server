package catalog

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/mande-server/action"
	"github.com/factly/mande-server/test"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestDeleteCatalog(t *testing.T) {

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

	t.Run("delete catalog", func(t *testing.T) {
		CatalogSelectMock(mock)

		productsAssociationSelectMock(mock, 1)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_catalog_product"`)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_catalog" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		e.DELETE(path).
			WithPath("catalog_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusOK)

		test.ExpectationsMet(t, mock)
	})

	t.Run("catalog record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(CatalogCols))

		e.DELETE(path).
			WithPath("catalog_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid catalog id", func(t *testing.T) {
		e.DELETE(path).
			WithPath("catalog_id", "abc").
			WithHeaders(headers).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("deleting catalog fails", func(t *testing.T) {
		CatalogSelectMock(mock)

		productsAssociationSelectMock(mock, 1)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_catalog_product"`)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_catalog" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnError(errors.New("cannot delete catalog"))
		mock.ExpectRollback()

		e.DELETE(path).
			WithPath("catalog_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("delete catalog when meili is down", func(t *testing.T) {
		gock.Off()
		test.KetoGock()
		test.KavachGock()
		gock.New(server.URL).EnableNetworking().Persist()

		CatalogSelectMock(mock)

		productsAssociationSelectMock(mock, 1)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "dp_catalog_product"`)).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_catalog" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectRollback()

		e.DELETE(path).
			WithPath("catalog_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
