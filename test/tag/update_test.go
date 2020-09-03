package tag

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestUpdateTag(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	test.MeiliGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	e := httpexpect.New(t, server.URL)

	t.Run("update tag", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(TagCols).
				AddRow(1, time.Now(), time.Now(), nil, "Original Tag", "original-tag"))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE \"dp_tag\" SET (.+)  WHERE (.+) \"dp_tag\".\"id\" = `).
			WithArgs(Tag["slug"], Tag["title"], test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		TagSelectMock(mock)
		mock.ExpectCommit()

		e.PUT(path).
			WithPath("tag_id", "1").
			WithJSON(Tag).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(Tag)

		test.ExpectationsMet(t, mock)
	})

	t.Run("tag not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(TagCols))

		e.PUT(path).
			WithPath("tag_id", "1").
			WithJSON(Tag).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid tag id", func(t *testing.T) {
		e.PUT(path).
			WithPath("tag_id", "abc").
			WithJSON(Tag).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("unprocessable tag body", func(t *testing.T) {
		e.PUT(path).
			WithPath("tag_id", "1").
			WithJSON(invalidTag).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("update tag when meili is down", func(t *testing.T) {
		gock.Off()
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(TagCols).
				AddRow(1, time.Now(), time.Now(), nil, "Original Tag", "original-tag"))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE \"dp_tag\" SET (.+)  WHERE (.+) \"dp_tag\".\"id\" = `).
			WithArgs(Tag["slug"], Tag["title"], test.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		TagSelectMock(mock)
		mock.ExpectRollback()

		e.PUT(path).
			WithPath("tag_id", "1").
			WithJSON(Tag).
			Expect().
			Status(http.StatusInternalServerError)
		test.ExpectationsMet(t, mock)
	})

}
