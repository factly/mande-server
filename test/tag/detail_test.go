package tag

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/mande-server/action"
	"github.com/factly/mande-server/test"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestDetailTag(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	adminExpect := httpexpect.New(t, server.URL)

	test.KavachGock()
	test.KetoGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	CommonDetailsTests(t, mock, adminExpect)

	server.Close()

	router = action.RegisterUserRoutes()
	server = httptest.NewServer(router)
	userExpect := httpexpect.New(t, server.URL)

	gock.New(server.URL).EnableNetworking().Persist()

	CommonDetailsTests(t, mock, userExpect)

	server.Close()
}

func CommonDetailsTests(t *testing.T, mock sqlmock.Sqlmock, e *httpexpect.Expect) {
	t.Run("get tag by id", func(t *testing.T) {
		TagSelectMock(mock, 1)

		e.GET(path).
			WithPath("tag_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			Keys().
			Contains("id", "created_at", "updated_at", "deleted_at", "title", "slug")

		test.ExpectationsMet(t, mock)
	})

	t.Run("tag record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(TagCols))

		e.GET(path).
			WithPath("tag_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid tag id", func(t *testing.T) {
		e.GET(path).
			WithPath("tag_id", "abc").
			WithHeaders(headers).
			Expect().
			Status(http.StatusBadRequest)
	})
}
