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

func TestGetTagDetail(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	// DB
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tag"`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "slug"}).
			AddRow(1, time.Now(), time.Now(), nil, "Test Tag", "test-tag"))

	// Request
	e.GET("/tags/1").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		Keys().
		Contains("id", "created_at", "updated_at", "deleted_at", "title", "slug")

	mock.ExpectationsWereMet()
}
