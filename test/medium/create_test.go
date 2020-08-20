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

func TestCreateMedium(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("create medium", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_medium"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, medium["name"], medium["slug"], medium["type"], medium["title"], medium["description"], medium["caption"], medium["alt_text"], medium["file_size"], medium["url"], medium["dimensions"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
		mock.ExpectCommit()

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(mediumCols).
				AddRow(1, time.Now(), time.Now(), nil, medium["name"], medium["slug"], medium["type"], medium["title"], medium["description"], medium["caption"], medium["alt_text"], medium["file_size"], medium["url"], medium["dimensions"]))

		e.POST(path).
			WithJSON(medium).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(medium)

		mock.ExpectationsWereMet()
	})

	t.Run("unprocessable medium body", func(t *testing.T) {
		e.POST(path).
			WithJSON(invalidMedium).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("empty medium body", func(t *testing.T) {
		e.POST(path).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})
}
