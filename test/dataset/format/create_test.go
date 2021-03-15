package datasetformat

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/mande-server/action"
	"github.com/factly/mande-server/test"
	"github.com/factly/mande-server/test/format"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestCreateDatasetFormat(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	test.KetoGock()
	test.KavachGock()
	gock.New(server.URL).EnableNetworking().Persist()

	e := httpexpect.New(t, server.URL)

	t.Run("create a dataset format", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_dataset_format"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, DatasetFormat["format_id"], 1, DatasetFormat["url"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		DatasetFormatSelectMock(mock)

		format.FormatSelectMock(mock)

		e.POST(basePath).
			WithPathObject(map[string]interface{}{
				"dataset_id": 1,
			}).
			WithJSON(DatasetFormat).
			WithHeaders(headers).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(DatasetFormat).
			Value("format").
			Object().
			ContainsMap(format.Format)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable dataset format body", func(t *testing.T) {
		e.POST(basePath).
			WithPathObject(map[string]interface{}{
				"dataset_id": 1,
			}).
			WithHeaders(headers).
			WithJSON(invalidDatasetFormat).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("undecodable dataset format body", func(t *testing.T) {
		e.POST(basePath).
			WithPathObject(map[string]interface{}{
				"dataset_id": 1,
			}).
			WithHeaders(headers).
			WithJSON(undecodableDatasetFormat).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("invalid dataset id", func(t *testing.T) {
		e.POST(basePath).
			WithPathObject(map[string]interface{}{
				"dataset_id": "abc",
			}).
			WithHeaders(headers).
			WithJSON(DatasetFormat).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("format does not exist", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_dataset_format"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, DatasetFormat["format_id"], 1, DatasetFormat["url"]).
			WillReturnError(errDatasetFormatFK)
		mock.ExpectRollback()

		e.POST(basePath).
			WithPath("dataset_id", "1").
			WithHeaders(headers).
			WithJSON(DatasetFormat).
			Expect().
			Status(http.StatusInternalServerError)
	})
}
