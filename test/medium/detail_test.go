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

func TestDetailMedium(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("get medium by id", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(mediumCols).
				AddRow(1, time.Now(), time.Now(), nil, medium["name"], medium["slug"], medium["type"], medium["title"], medium["description"], medium["caption"], medium["alt_text"], medium["file_size"], medium["url"], medium["dimensions"]))

		e.GET(pathId).
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

		e.GET(pathId).
			Expect().
			Status(http.StatusNotFound)

		mock.ExpectationsWereMet()
	})

	t.Run("invalid id", func(t *testing.T) {
		e.GET(pathInvalidId).
			Expect().
			Status(http.StatusNotFound)
	})
}
