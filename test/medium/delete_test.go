package medium

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/gavv/httpexpect"
)

func TestDeleteMedium(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("delete medium", func(t *testing.T) {
		mediumSelectMock(mock)

		mediumCatalogExpect(mock, 0)

		mediumDatasetExpect(mock, 0)

		mediumProductExpect(mock, 0)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_medium" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		e.DELETE(path).
			WithPath("media_id", "1").
			Expect().
			Status(http.StatusOK)

		test.ExpectationsMet(t, mock)
	})

	t.Run("medium record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(mediumCols))

		e.DELETE(path).
			WithPath("media_id", "1").
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid medium id", func(t *testing.T) {
		e.DELETE(path).
			WithPath("media_id", "abc").
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("medium is associated with catalog", func(t *testing.T) {
		mediumSelectMock(mock)

		mediumCatalogExpect(mock, 1)

		e.DELETE(path).
			WithPath("media_id", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})

	t.Run("medium is associated with dataset", func(t *testing.T) {
		mediumSelectMock(mock)

		mediumCatalogExpect(mock, 0)

		mediumDatasetExpect(mock, 1)

		e.DELETE(path).
			WithPath("media_id", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})

	t.Run("medium is associated with product", func(t *testing.T) {
		mediumSelectMock(mock)

		mediumCatalogExpect(mock, 0)

		mediumDatasetExpect(mock, 0)

		mediumProductExpect(mock, 1)

		e.DELETE(path).
			WithPath("media_id", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})
}
