package tag

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/h2non/gock.v1"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/gavv/httpexpect"
)

func TestCreateTag(t *testing.T) {
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

	t.Run("create a tag", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_tag"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, Tag["title"], Tag["slug"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		e.POST(basePath).
			WithJSON(Tag).
			WithHeaders(headers).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(Tag)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable tag body", func(t *testing.T) {
		e.POST(basePath).
			WithJSON(invalidTag).
			WithHeaders(headers).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("empty tag body", func(t *testing.T) {
		e.POST(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("creating tag fails", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_tag"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, Tag["title"], Tag["slug"]).
			WillReturnError(errors.New("cannot create"))
		mock.ExpectRollback()

		e.POST(basePath).
			WithJSON(Tag).
			WithHeaders(headers).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("create a tag when meili is down", func(t *testing.T) {
		gock.Off()
		test.KavachGock()
		test.KetoGock()
		gock.New(server.URL).EnableNetworking().Persist()

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_tag"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, Tag["title"], Tag["slug"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectRollback()

		e.POST(basePath).
			WithJSON(Tag).
			WithHeaders(headers).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
