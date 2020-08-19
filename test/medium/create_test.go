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

func TestCreateMedium(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	createdMedium := map[string]interface{}{
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

	t.Run("create medium", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_medium"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, createdMedium["name"], createdMedium["slug"], createdMedium["type"], createdMedium["title"], createdMedium["description"], createdMedium["caption"], createdMedium["alt_text"], createdMedium["file_size"], createdMedium["url"], createdMedium["dimensions"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
		mock.ExpectCommit()

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(mediumCols).
				AddRow(1, time.Now(), time.Now(), nil, createdMedium["name"], createdMedium["slug"], createdMedium["type"], createdMedium["title"], createdMedium["description"], createdMedium["caption"], createdMedium["alt_text"], createdMedium["file_size"], createdMedium["url"], createdMedium["dimensions"]))

		e.POST("/media").
			WithJSON(createdMedium).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(createdMedium)

		mock.ExpectationsWereMet()
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

		e.POST("/media").
			WithJSON(invalidMedium).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("empty medium body", func(t *testing.T) {
		e.POST("/media").
			Expect().
			Status(http.StatusUnprocessableEntity)
	})
}
