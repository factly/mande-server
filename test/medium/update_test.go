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
	test.KavachGock()
	test.KetoGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	e := httpexpect.New(t, server.URL)

	t.Run("update medium", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(MediumCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, "name", "slug", "type", "title", "description", "caption", "alt_text", 100, nil, "dimensions"))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE \"dp_medium\"`).
			WithArgs(test.AnyTime{}, 1, Medium["name"], Medium["slug"], Medium["type"], Medium["title"], Medium["description"], Medium["caption"], Medium["alt_text"], Medium["file_size"], Medium["url"], Medium["dimensions"], 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		MediumSelectMock(mock, 1, 1)
		mock.ExpectCommit()

		e.PUT(path).
			WithPath("media_id", "1").
			WithHeaders(headers).
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
			WithHeaders(headers).
			WithJSON(Medium).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid medium id", func(t *testing.T) {
		e.PUT(path).
			WithPath("media_id", "abc").
			WithHeaders(headers).
			WithJSON(Medium).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("unprocessable medium body", func(t *testing.T) {
		e.PUT(path).
			WithPath("media_id", "1").
			WithHeaders(headers).
			WithJSON(invalidMedium).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("undecodable medium body", func(t *testing.T) {
		e.PUT(path).
			WithPath("media_id", "1").
			WithHeaders(headers).
			WithJSON(undecodableMedium).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("update medium when meili is down", func(t *testing.T) {
		gock.Off()
		test.KavachGock()
		test.KetoGock()
		gock.New(server.URL).EnableNetworking().Persist()

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(MediumCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, "name", "slug", "type", "title", "description", "caption", "alt_text", 100, nil, "dimensions"))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE \"dp_medium\"`).
			WithArgs(test.AnyTime{}, 1, Medium["name"], Medium["slug"], Medium["type"], Medium["title"], Medium["description"], Medium["caption"], Medium["alt_text"], Medium["file_size"], Medium["url"], Medium["dimensions"], 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		MediumSelectMock(mock, 1, 1)
		mock.ExpectRollback()

		e.PUT(path).
			WithPath("media_id", "1").
			WithHeaders(headers).
			WithJSON(Medium).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
