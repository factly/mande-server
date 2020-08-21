package format

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/gavv/httpexpect"
)

func TestCreateFormat(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("create a format", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_format"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, format["name"], format["description"], format["is_default"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		formatSelectMock(mock)

		e.POST(basePath).
			WithJSON(format).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(format)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable format body", func(t *testing.T) {
		e.POST(basePath).
			WithJSON(invalidFormat).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("empty format body", func(t *testing.T) {
		e.POST(basePath).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})
}
