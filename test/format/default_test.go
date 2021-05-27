package format

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/mande-server/action"
	"github.com/factly/mande-server/action/format"
	"github.com/factly/mande-server/test"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestCreateDefaultFormat(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	format.DataFile = "../../data/format.json"

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	test.KetoGock()
	test.KavachGock()
	gock.New(server.URL).EnableNetworking().Persist()

	t.Run("create default format", func(t *testing.T) {
		mock.ExpectBegin()
		for i := 0; i < 3; i++ {
			mock.ExpectQuery(selectQuery).
				WillReturnRows(sqlmock.NewRows(FormatCols))

			mock.ExpectQuery(`INSERT INTO "dp_format"`).
				WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, defaultFormats[i]["name"], defaultFormats[i]["description"], defaultFormats[i]["is_default"]).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		}
		mock.ExpectCommit()

		e.POST(defaultPath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusCreated).JSON().
			Object().
			Value("nodes").
			Array().Element(0).Object().ContainsMap(defaultFormats[0])
		test.ExpectationsMet(t, mock)
	})

	t.Run("default formats already created", func(t *testing.T) {
		mock.ExpectBegin()
		for i := 0; i < 3; i++ {
			mock.ExpectQuery(selectQuery).
				WillReturnRows(sqlmock.NewRows(FormatCols).
					AddRow(1, time.Now(), time.Now(), nil, 1, 1, defaultFormats[i]["name"], defaultFormats[i]["description"], defaultFormats[i]["is_default"]))
		}

		mock.ExpectCommit()

		e.POST(defaultPath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusCreated).JSON().
			Object().
			Value("nodes").
			Array().Element(0).Object().ContainsMap(defaultFormats[0])
		test.ExpectationsMet(t, mock)
	})

	t.Run("cannot open data file", func(t *testing.T) {
		format.DataFile = "nofile.json"
		e.POST(defaultPath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusInternalServerError)
		format.DataFile = "../../data/format.json"
	})

	t.Run("cannot parse data file", func(t *testing.T) {
		format.DataFile = "invalidData.json"
		e.POST(defaultPath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusInternalServerError)
		format.DataFile = "../../data/format.json"
	})
}
