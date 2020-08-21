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

func TestDetailTag(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("get tag by id", func(t *testing.T) {
		TagSelectMock(mock)

		e.GET(path).
			WithPath("tag_id", "1").
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
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid tag id", func(t *testing.T) {
		e.GET(path).
			WithPath("tag_id", "abc").
			Expect().
			Status(http.StatusNotFound)
	})

}
