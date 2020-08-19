package tag

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

func TestUpdateTag(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	updatedTag := map[string]interface{}{
		"title": "Test Updated Tag",
		"slug":  "test-updated-tag",
	}
	selectQuery := regexp.QuoteMeta(`SELECT * FROM "dp_tag"`)

	t.Run("update tag", func(t *testing.T) {

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "slug"}).
				AddRow(1, time.Now(), time.Now(), nil, "Original Tag", "original-tag"))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE \"dp_tag\" SET (.+)  WHERE (.+) \"dp_tag\".\"id\" = `).
			WithArgs(updatedTag["slug"], updatedTag["title"], test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "slug"}).
				AddRow(1, time.Now(), time.Now(), nil, updatedTag["title"], updatedTag["slug"]))

		e.PUT("/tags/1").
			WithJSON(updatedTag).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(updatedTag)

		mock.ExpectationsWereMet()
	})

	t.Run("tag not found", func(t *testing.T) {

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "slug"}))

		e.PUT("/tags/1").
			WithJSON(updatedTag).
			Expect().
			Status(http.StatusNotFound)

		mock.ExpectationsWereMet()
	})

	t.Run("invalid tag id", func(t *testing.T) {

		e.PUT("/tags/abc").
			WithJSON(updatedTag).
			Expect().
			Status(http.StatusNotFound)

	})

	t.Run("unprocessable tag body", func(t *testing.T) {

		invalidTag := map[string]interface{}{
			"titl": "Test",
			"slg":  "test",
		}

		e.POST("/tags").
			WithJSON(invalidTag).
			Expect().
			Status(http.StatusUnprocessableEntity)

	})

}
