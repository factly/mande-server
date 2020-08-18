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

	//CREATE TAG
	updatedTag := map[string]interface{}{
		"title": "Test Updated Tag",
		"slug":  "test-updated-tag",
	}
	tagID := 1

	// DB
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE \"dp_tag\" SET (.+)  WHERE (.+) \"dp_tag\".\"id\" = `).
		WithArgs(updatedTag["slug"], updatedTag["title"], test.AnyTime{}, tagID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tag"`)).
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

}
