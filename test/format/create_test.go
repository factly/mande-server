package format

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestCreateFormat(t *testing.T) {
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

	t.Run("create a format", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_format"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, Format["name"], Format["description"], Format["is_default"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		e.POST(basePath).
			WithJSON(Format).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(Format)

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

	t.Run("creating format fails", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_format"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, Format["name"], Format["description"], Format["is_default"]).
			WillReturnError(errors.New("cannot create format"))
		mock.ExpectRollback()

		e.POST(basePath).
			WithJSON(Format).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("create a format when meili is down", func(t *testing.T) {
		gock.Off()
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_format"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, Format["name"], Format["description"], Format["is_default"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectRollback()

		e.POST(basePath).
			WithJSON(Format).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
