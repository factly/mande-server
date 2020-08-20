package medium

import (
	"net/http"
	"net/http/httptest"
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

	t.Run("update medium", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(mediumCols).
				AddRow(1, time.Now(), time.Now(), nil, "name", "slug", "type", "title", "description", "caption", "alt_text", 100, "url", "dimensions"))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE \"dp_medium\" SET (.+)  WHERE (.+) \"dp_medium\".\"id\" = `).
			WithArgs(medium["alt_text"], medium["caption"], medium["description"], medium["dimensions"], medium["file_size"], medium["name"], medium["slug"], medium["title"], medium["type"], test.AnyTime{}, medium["url"], 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(mediumCols).
				AddRow(1, time.Now(), time.Now(), nil, medium["name"], medium["slug"], medium["type"], medium["title"], medium["description"], medium["caption"], medium["alt_text"], medium["file_size"], medium["url"], medium["dimensions"]))

		e.PUT(pathId).
			WithJSON(medium).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(medium)

		mock.ExpectationsWereMet()
	})

	t.Run("medium record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(mediumCols))

		e.PUT(pathId).
			WithJSON(medium).
			Expect().
			Status(http.StatusNotFound)

		mock.ExpectationsWereMet()
	})

	t.Run("invalid medium id", func(t *testing.T) {
		e.PUT(pathInvalidId).
			WithJSON(medium).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("unprocessable medium body", func(t *testing.T) {
		e.PUT(pathId).
			WithJSON(invalidMedium).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

}
