package datasetformat

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/mande-server/action"
	"github.com/factly/mande-server/test"
	"github.com/factly/mande-server/test/format"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestListDatasetFormat(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	test.KetoGock()
	test.KavachGock()
	gock.New(server.URL).EnableNetworking().Persist()

	t.Run("get empty list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(DatasetFormatCols))

		e.GET(basePath).
			WithHeaders(headers).
			WithPath("dataset_id", "1").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 0})

		test.ExpectationsMet(t, mock)
	})

	t.Run("list dataset formats", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(datasetformatlist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(DatasetFormatCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, datasetformatlist[0]["format_id"], datasetformatlist[0]["dataset_id"], datasetformatlist[0]["url"]).
				AddRow(2, time.Now(), time.Now(), nil, 1, 1, datasetformatlist[1]["format_id"], datasetformatlist[1]["dataset_id"], datasetformatlist[1]["url"]))

		format.FormatSelectMock(mock)

		e.GET(basePath).
			WithHeaders(headers).
			WithPath("dataset_id", "1").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(datasetformatlist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(datasetformatlist[0])

		test.ExpectationsMet(t, mock)
	})

	t.Run("list dataset formats with paiganation", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(datasetformatlist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(DatasetFormatCols).
				AddRow(2, time.Now(), time.Now(), nil, 1, 1, datasetformatlist[1]["format_id"], datasetformatlist[1]["dataset_id"], datasetformatlist[1]["url"]))

		format.FormatSelectMock(mock)

		e.GET(basePath).
			WithHeaders(headers).
			WithPath("dataset_id", "1").
			WithQueryObject(map[string]interface{}{
				"limit": "1",
				"page":  "2",
			}).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(datasetformatlist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(datasetformatlist[1])

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid dataset id", func(t *testing.T) {
		e.GET(basePath).
			WithHeaders(headers).
			WithPath("dataset_id", "abc").
			Expect().
			Status(http.StatusBadRequest)
	})
}
