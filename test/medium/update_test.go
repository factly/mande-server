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

func TestUpdateMedium(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	updatedMedium := map[string]interface{}{
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

	t.Run("update medium", func(t *testing.T) {

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(mediumCols).
				AddRow(1, time.Now(), time.Now(), nil, "name", "slug", "type", "title", "description", "caption", "alt_text", 100, "url", "dimensions"))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE \"dp_medium\" SET (.+)  WHERE (.+) \"dp_medium\".\"id\" = `).
			WithArgs(updatedMedium["alt_text"], updatedMedium["caption"], updatedMedium["description"], updatedMedium["dimensions"], updatedMedium["file_size"], updatedMedium["name"], updatedMedium["slug"], updatedMedium["title"], updatedMedium["type"], test.AnyTime{}, updatedMedium["url"], 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(mediumCols).
				AddRow(1, time.Now(), time.Now(), nil, updatedMedium["name"], updatedMedium["slug"], updatedMedium["type"], updatedMedium["title"], updatedMedium["description"], updatedMedium["caption"], updatedMedium["alt_text"], updatedMedium["file_size"], updatedMedium["url"], updatedMedium["dimensions"]))

		e.PUT("/media/1").
			WithJSON(updatedMedium).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(updatedMedium)

		mock.ExpectationsWereMet()

	})

	t.Run("medium record not found", func(t *testing.T) {

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(mediumCols))

		e.PUT("/media/1").
			WithJSON(updatedMedium).
			Expect().
			Status(http.StatusNotFound)

		mock.ExpectationsWereMet()
	})

	t.Run("invalid medium id", func(t *testing.T) {
		e.PUT("/media/abc").
			WithJSON(updatedMedium).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("unprocessable medium body", func(t *testing.T) {
		invalidMedium := map[string]interface{}{
			"nam":         "Test Medium",
			"slug":        "test-medium",
			"type":        "testtype",
			"title":       "Test Title",
			"description": "Test Description",
			"caption":     "Test Caption",
			"alt_text":    "Test alt text",
			"filesize":    100,
			"url":         "http:/testurl.com",
			"dimensions":  "testdims",
		}

		e.PUT("/media/1").
			WithJSON(invalidMedium).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

}
