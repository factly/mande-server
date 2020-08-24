package datasetformat

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/format"
	"github.com/gavv/httpexpect"
)

func TestCreateDatasetFormat(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("create a dataset format", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_dataset_format"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, DatasetFormat["format_id"], 1, DatasetFormat["url"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		DatasetFormatSelectMock(mock)

		format.FormatSelectMock(mock)

		e.POST(basePath).
			WithPathObject(map[string]interface{}{
				"dataset_id": 1,
			}).
			WithJSON(DatasetFormat).
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
			WithJSON(invalidDatasetFormat).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("format does not exist", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_dataset_format"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, DatasetFormat["format_id"], 1, DatasetFormat["url"]).
			WillReturnError(errDatasetFormatFK)
		mock.ExpectRollback()

		e.POST(basePath).
			WithPath("dataset_id", "1").
			WithJSON(DatasetFormat).
			Expect().
			Status(http.StatusInternalServerError)
	})
}
