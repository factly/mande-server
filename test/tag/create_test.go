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

func TestCreateTag(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	//CREATE TAG
	createdTag := map[string]interface{}{
		"title": "Test Tag",
		"slug":  "test-tag",
	}

	tagCols := []string{"id", "created_at", "updated_at", "deleted_at", "title", "slug"}

	t.Run("create a tag", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_tag"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, createdTag["title"], createdTag["slug"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tag"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(tagCols).
				AddRow(1, time.Now(), time.Now(), nil, createdTag["title"], createdTag["slug"]))

		e.POST("/tags").
			WithJSON(createdTag).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(createdTag)

		mock.ExpectationsWereMet()
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

	t.Run("empty tag body", func(t *testing.T) {

		e.POST("/tags").
			Expect().
			Status(http.StatusUnprocessableEntity)

	})

}
