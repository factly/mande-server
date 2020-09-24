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
	"gopkg.in/h2non/gock.v1"
)

func TestUpdateMedium(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	test.MeiliGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	e := httpexpect.New(t, server.URL)

	t.Run("update medium", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(MediumCols).
				AddRow(1, time.Now(), time.Now(), nil, "name", "slug", "type", "title", "description", "caption", "alt_text", 100, "url", "dimensions"))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE \"dp_medium\" SET (.+)  WHERE (.+) \"dp_medium\".\"id\" = `).
			WithArgs(Medium["alt_text"], Medium["caption"], Medium["description"], Medium["dimensions"], Medium["file_size"], Medium["name"], Medium["slug"], Medium["title"], Medium["type"], test.AnyTime{}, Medium["url"], 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		MediumSelectMock(mock)
		mock.ExpectCommit()

		e.PUT(path).
			WithPath("media_id", "1").
			WithJSON(Medium).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(Medium)

		test.ExpectationsMet(t, mock)
	})

	t.Run("medium record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(MediumCols))

		e.PUT(path).
			WithPath("media_id", "1").
			WithJSON(Medium).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid medium id", func(t *testing.T) {
		e.PUT(path).
			WithPath("media_id", "abc").
			WithJSON(Medium).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("unprocessable medium body", func(t *testing.T) {
		e.PUT(path).
			WithPath("media_id", "1").
			WithJSON(invalidMedium).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("undecodable medium body", func(t *testing.T) {
		e.PUT(path).
			WithPath("media_id", "1").
			WithJSON(undecodableMedium).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("update medium when meili is down", func(t *testing.T) {
		gock.Off()
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(MediumCols).
				AddRow(1, time.Now(), time.Now(), nil, "name", "slug", "type", "title", "description", "caption", "alt_text", 100, "url", "dimensions"))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE \"dp_medium\" SET (.+)  WHERE (.+) \"dp_medium\".\"id\" = `).
			WithArgs(Medium["alt_text"], Medium["caption"], Medium["description"], Medium["dimensions"], Medium["file_size"], Medium["name"], Medium["slug"], Medium["title"], Medium["type"], test.AnyTime{}, Medium["url"], 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		MediumSelectMock(mock)
		mock.ExpectRollback()

		e.PUT(path).
			WithPath("media_id", "1").
			WithJSON(Medium).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
