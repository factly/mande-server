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

func TestDetailTag(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	tagCols := []string{"id", "created_at", "updated_at", "deleted_at", "title", "slug"}
	selectQuery := regexp.QuoteMeta(`SELECT * FROM "dp_tag"`)

	t.Run("get tag by id", func(t *testing.T) {

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(tagCols).
				AddRow(1, time.Now(), time.Now(), nil, "Test Tag", "test-tag"))

		e.GET("/tags/1").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			Keys().
			Contains("id", "created_at", "updated_at", "deleted_at", "title", "slug")

		mock.ExpectationsWereMet()
	})

	t.Run("tag record not found", func(t *testing.T) {

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(tagCols))

		e.GET("/tags/1").
			Expect().
			Status(http.StatusNotFound)

		mock.ExpectationsWereMet()
	})

	t.Run("invalid tag id", func(t *testing.T) {

		e.GET("/tags/abc").
			Expect().
			Status(http.StatusNotFound)

	})

}
