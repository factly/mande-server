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

func TestDeleteTag(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	// Test objects
	deletedTag := map[string]interface{}{
		"title": "Test Delete Tag",
		"slug":  "test-delete-tag",
	}

	selectQuery := regexp.QuoteMeta(`SELECT * FROM "dp_tag"`)
	tagProductQuery := regexp.QuoteMeta(`SELECT count(*) FROM "dp_product" INNER JOIN "dp_product_tag"`)
	tagDatasetQuery := regexp.QuoteMeta(`SELECT count(*) FROM "dp_dataset" INNER JOIN "dp_dataset_tag"`)

	t.Run("delete tag", func(t *testing.T) {

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "slug"}).
				AddRow(1, time.Now(), time.Now(), nil, deletedTag["title"], deletedTag["slug"]))

		mock.ExpectQuery(tagProductQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(tagDatasetQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_tag" SET "deleted_at"=`)).
			WithArgs(test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		e.DELETE("/tags/1").
			Expect().
			Status(http.StatusOK)

		mock.ExpectationsWereMet()

	})

	t.Run("tag not found", func(t *testing.T) {

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "slug"}))

		e.DELETE("/tags/1").
			Expect().
			Status(http.StatusNotFound)

		mock.ExpectationsWereMet()
	})

	t.Run("invalid tag id", func(t *testing.T) {

		e.DELETE("/tags/abc").
			Expect().
			Status(http.StatusNotFound)

	})

	t.Run("tag associated with product", func(t *testing.T) {

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "slug"}).
				AddRow(1, time.Now(), time.Now(), nil, deletedTag["title"], deletedTag["slug"]))

		mock.ExpectQuery(tagProductQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("1"))

		e.DELETE("/tags/1").
			Expect().
			Status(http.StatusUnprocessableEntity)

		mock.ExpectationsWereMet()

	})

	t.Run("tag associated with dataset", func(t *testing.T) {

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "title", "slug"}).
				AddRow(1, time.Now(), time.Now(), nil, deletedTag["title"], deletedTag["slug"]))

		mock.ExpectQuery(tagProductQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(tagDatasetQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("1"))

		e.DELETE("/tags/1").
			Expect().
			Status(http.StatusUnprocessableEntity)

		mock.ExpectationsWereMet()

	})
}
