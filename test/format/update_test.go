package format

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

func TestUpdateFormat(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("update format", func(t *testing.T) {

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(formatCols).
				AddRow(1, time.Now(), time.Now(), nil, "name", "description", true))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE \"dp_format\" SET (.+)  WHERE (.+) \"dp_format\".\"id\" = `).
			WithArgs(format["description"], format["is_default"], format["name"], test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		formatSelectMock(mock)

		e.PUT(path).
			WithPath("format_id", "1").
			WithJSON(format).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(format)

		test.ExpectationsMet(t, mock)
	})

	t.Run("format record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(formatCols))

		e.PUT(path).
			WithPath("format_id", "1").
			WithJSON(format).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable format body", func(t *testing.T) {
		e.PUT(path).
			WithPath("format_id", "1").
			WithJSON(invalidFormat).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("invalid format id", func(t *testing.T) {
		e.PUT(path).
			WithPath("format_id", "abc").
			WithJSON(format).
			Expect().
			Status(http.StatusNotFound)
	})
}
