package tag

import (
	"net/http"
	"net/http/httptest"
	"testing"

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

	t.Run("create a tag", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_tag"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, tag["title"], tag["slug"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		tagSelectMock(mock)

		e.POST(basePath).
			WithJSON(tag).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(tag)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable tag body", func(t *testing.T) {
		e.POST(basePath).
			WithJSON(invalidTag).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("empty tag body", func(t *testing.T) {
		e.POST(basePath).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

}
