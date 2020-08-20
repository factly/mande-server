package medium

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

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

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(mediumCols).
				AddRow(1, time.Now(), time.Now(), nil, medium["name"], medium["slug"], medium["type"], medium["title"], medium["description"], medium["caption"], medium["alt_text"], medium["file_size"], medium["url"], medium["dimensions"]))

		mock.ExpectQuery(mediumCatalogQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(mediumDatasetQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(mediumProductQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_medium" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		e.DELETE(pathId).
			Expect().
			Status(http.StatusOK)

		mock.ExpectationsWereMet()
	})

	t.Run("medium record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(mediumCols))

		e.DELETE(pathId).
			Expect().
			Status(http.StatusNotFound)

		mock.ExpectationsWereMet()
	})

	t.Run("invalid medium id", func(t *testing.T) {
		e.DELETE(pathInvalidId).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("medium is associated with catalog", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(mediumCols).
				AddRow(1, time.Now(), time.Now(), nil, medium["name"], medium["slug"], medium["type"], medium["title"], medium["description"], medium["caption"], medium["alt_text"], medium["file_size"], medium["url"], medium["dimensions"]))

		mock.ExpectQuery(mediumCatalogQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("1"))

		e.DELETE(pathId).
			Expect().
			Status(http.StatusUnprocessableEntity)

		mock.ExpectationsWereMet()
	})

	t.Run("medium is associated with dataset", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(mediumCols).
				AddRow(1, time.Now(), time.Now(), nil, medium["name"], medium["slug"], medium["type"], medium["title"], medium["description"], medium["caption"], medium["alt_text"], medium["file_size"], medium["url"], medium["dimensions"]))

		mock.ExpectQuery(mediumCatalogQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(mediumDatasetQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("1"))

		e.DELETE(pathId).
			Expect().
			Status(http.StatusUnprocessableEntity)

		mock.ExpectationsWereMet()
	})

	t.Run("medium is associated with product", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(mediumCols).
				AddRow(1, time.Now(), time.Now(), nil, medium["name"], medium["slug"], medium["type"], medium["title"], medium["description"], medium["caption"], medium["alt_text"], medium["file_size"], medium["url"], medium["dimensions"]))

		mock.ExpectQuery(mediumCatalogQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(mediumDatasetQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(mediumProductQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("1"))

		e.DELETE(pathId).
			Expect().
			Status(http.StatusUnprocessableEntity)

		mock.ExpectationsWereMet()
	})
}
