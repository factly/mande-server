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

	deletedMedium := map[string]interface{}{
		"name":        "Test Medium",
		"slug":        "test-medium",
		"type":        "testtype",
		"title":       "Test Title",
		"description": "Test Description",
		"caption":     "Test Caption",
		"alt_text":    "Test alt text",
		"file_size":   100,
		"url":         "http:/testurl.com",
		"dimensions":  "testdims",
	}
	mediumCols := []string{"id", "created_at", "updated_at", "deleted_at", "name", "slug", "type", "title", "description", "caption", "alt_text", "file_size", "url", "dimensions"}
	selectQuery := regexp.QuoteMeta(`SELECT * FROM "dp_medium"`)
	mediumCatalogQuery := regexp.QuoteMeta(`SELECT count(*) FROM "dp_catalog`)
	mediumDatasetQuery := regexp.QuoteMeta(`SELECT count(*) FROM "dp_dataset`)
	mediumProductQuery := regexp.QuoteMeta(`SELECT count(*) FROM "dp_product`)

	t.Run("delete medium", func(t *testing.T) {

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(mediumCols).
				AddRow(1, time.Now(), time.Now(), nil, deletedMedium["name"], deletedMedium["slug"], deletedMedium["type"], deletedMedium["title"], deletedMedium["description"], deletedMedium["caption"], deletedMedium["alt_text"], deletedMedium["file_size"], deletedMedium["url"], deletedMedium["dimensions"]))

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

		e.DELETE("/media/1").
			Expect().
			Status(http.StatusOK)

		mock.ExpectationsWereMet()
	})

	t.Run("medium record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(mediumCols))

		e.DELETE("/media/1").
			Expect().
			Status(http.StatusNotFound)

		mock.ExpectationsWereMet()
	})

	t.Run("invalid medium id", func(t *testing.T) {
		e.DELETE("/media/abc").
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("medium is associated with catalog", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(mediumCols).
				AddRow(1, time.Now(), time.Now(), nil, deletedMedium["name"], deletedMedium["slug"], deletedMedium["type"], deletedMedium["title"], deletedMedium["description"], deletedMedium["caption"], deletedMedium["alt_text"], deletedMedium["file_size"], deletedMedium["url"], deletedMedium["dimensions"]))

		mock.ExpectQuery(mediumCatalogQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("1"))

		e.DELETE("/media/1").
			Expect().
			Status(http.StatusUnprocessableEntity)

		mock.ExpectationsWereMet()
	})

	t.Run("medium is associated with dataset", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(mediumCols).
				AddRow(1, time.Now(), time.Now(), nil, deletedMedium["name"], deletedMedium["slug"], deletedMedium["type"], deletedMedium["title"], deletedMedium["description"], deletedMedium["caption"], deletedMedium["alt_text"], deletedMedium["file_size"], deletedMedium["url"], deletedMedium["dimensions"]))

		mock.ExpectQuery(mediumCatalogQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(mediumDatasetQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("1"))

		e.DELETE("/media/1").
			Expect().
			Status(http.StatusUnprocessableEntity)

		mock.ExpectationsWereMet()
	})

	t.Run("medium is associated with product", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(mediumCols).
				AddRow(1, time.Now(), time.Now(), nil, deletedMedium["name"], deletedMedium["slug"], deletedMedium["type"], deletedMedium["title"], deletedMedium["description"], deletedMedium["caption"], deletedMedium["alt_text"], deletedMedium["file_size"], deletedMedium["url"], deletedMedium["dimensions"]))

		mock.ExpectQuery(mediumCatalogQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(mediumDatasetQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(mediumProductQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("1"))

		e.DELETE("/media/1").
			Expect().
			Status(http.StatusUnprocessableEntity)

		mock.ExpectationsWereMet()
	})
}
