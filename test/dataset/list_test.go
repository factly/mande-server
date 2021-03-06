package dataset

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/medium"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestListDataset(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	adminExpect := httpexpect.New(t, server.URL)

	test.KetoGock()
	test.KavachGock()
	gock.New(server.URL).EnableNetworking().Persist()

	AdminListTests(t, mock, adminExpect)

	server.Close()

	router = action.RegisterUserRoutes()
	server = httptest.NewServer(router)
	userExpect := httpexpect.New(t, server.URL)

	gock.New(server.URL).EnableNetworking().Persist()

	UserListTests(t, mock, userExpect)

	server.Close()
}

func AdminListTests(t *testing.T, mock sqlmock.Sqlmock, e *httpexpect.Expect) {
	t.Run("get empty dataset list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(DatasetCols))

		e.GET(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 0})

		test.ExpectationsMet(t, mock)
	})

	t.Run("get dataset list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(datasetlist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(DatasetCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, datasetlist[0]["title"], datasetlist[0]["description"], datasetlist[0]["source"], datasetlist[0]["frequency"], datasetlist[0]["temporal_coverage"], datasetlist[0]["granularity"], datasetlist[0]["contact_name"], datasetlist[0]["contact_email"], datasetlist[0]["license"], datasetlist[0]["data_standard"], datasetlist[0]["sample_url"], datasetlist[0]["related_articles"], datasetlist[0]["time_saved"], datasetlist[0]["price"], datasetlist[0]["currency_id"], datasetlist[0]["featured_medium_id"]).
				AddRow(2, time.Now(), time.Now(), nil, 1, 1, datasetlist[1]["title"], datasetlist[1]["description"], datasetlist[1]["source"], datasetlist[1]["frequency"], datasetlist[1]["temporal_coverage"], datasetlist[1]["granularity"], datasetlist[1]["contact_name"], datasetlist[1]["contact_email"], datasetlist[1]["license"], datasetlist[1]["data_standard"], datasetlist[1]["sample_url"], datasetlist[1]["related_articles"], datasetlist[1]["time_saved"], datasetlist[1]["price"], datasetlist[1]["currency_id"], datasetlist[1]["featured_medium_id"]))

		currency.CurrencySelectMock(mock)
		medium.MediumSelectMock(mock)

		tagAssociationSelectMock(mock, 1, 2)

		datasetFormatSelectMock(mock, 1)

		datasetFormatSelectMock(mock, 2)

		delete(datasetlist[0], "tag_ids")

		e.GET(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(datasetlist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(datasetlist[0])

		test.ExpectationsMet(t, mock)
	})

	t.Run("get dataset list paiganation", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(datasetlist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(DatasetCols).
				AddRow(2, time.Now(), time.Now(), nil, 1, 1, datasetlist[1]["title"], datasetlist[1]["description"], datasetlist[1]["source"], datasetlist[1]["frequency"], datasetlist[1]["temporal_coverage"], datasetlist[1]["granularity"], datasetlist[1]["contact_name"], datasetlist[1]["contact_email"], datasetlist[1]["license"], datasetlist[1]["data_standard"], datasetlist[1]["sample_url"], datasetlist[1]["related_articles"], datasetlist[1]["time_saved"], datasetlist[1]["price"], datasetlist[1]["currency_id"], datasetlist[1]["featured_medium_id"]))

		currency.CurrencySelectMock(mock)
		medium.MediumSelectMock(mock)

		tagAssociationSelectMock(mock, 2)

		datasetFormatSelectMock(mock, 2)

		delete(datasetlist[1], "tag_ids")

		e.GET(basePath).
			WithQueryObject(map[string]interface{}{
				"limit": "1",
				"page":  "2",
			}).
			WithHeaders(headers).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(datasetlist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(datasetlist[1])

		test.ExpectationsMet(t, mock)
	})
}

func UserListTests(t *testing.T, mock sqlmock.Sqlmock, e *httpexpect.Expect) {
	t.Run("get empty dataset list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(DatasetCols))

		e.GET(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 0})

		test.ExpectationsMet(t, mock)
	})

	t.Run("get dataset list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(datasetlist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(DatasetCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, datasetlist[0]["title"], datasetlist[0]["description"], datasetlist[0]["source"], datasetlist[0]["frequency"], datasetlist[0]["temporal_coverage"], datasetlist[0]["granularity"], datasetlist[0]["contact_name"], datasetlist[0]["contact_email"], datasetlist[0]["license"], datasetlist[0]["data_standard"], datasetlist[0]["sample_url"], datasetlist[0]["related_articles"], datasetlist[0]["time_saved"], datasetlist[0]["price"], datasetlist[0]["currency_id"], datasetlist[0]["featured_medium_id"]).
				AddRow(2, time.Now(), time.Now(), nil, 1, 1, datasetlist[1]["title"], datasetlist[1]["description"], datasetlist[1]["source"], datasetlist[1]["frequency"], datasetlist[1]["temporal_coverage"], datasetlist[1]["granularity"], datasetlist[1]["contact_name"], datasetlist[1]["contact_email"], datasetlist[1]["license"], datasetlist[1]["data_standard"], datasetlist[1]["sample_url"], datasetlist[1]["related_articles"], datasetlist[1]["time_saved"], datasetlist[1]["price"], datasetlist[1]["currency_id"], datasetlist[1]["featured_medium_id"]))

		currency.CurrencySelectMock(mock)
		medium.MediumSelectMock(mock)

		tagAssociationSelectMock(mock, 1, 2)

		userDatasetFormatSelectMock(mock, 1)

		userDatasetFormatSelectMock(mock, 2)

		delete(datasetlist[0], "tag_ids")

		e.GET(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(datasetlist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(datasetlist[0])

		test.ExpectationsMet(t, mock)
	})

	t.Run("get dataset list paiganation", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(datasetlist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(DatasetCols).
				AddRow(2, time.Now(), time.Now(), nil, 1, 1, datasetlist[1]["title"], datasetlist[1]["description"], datasetlist[1]["source"], datasetlist[1]["frequency"], datasetlist[1]["temporal_coverage"], datasetlist[1]["granularity"], datasetlist[1]["contact_name"], datasetlist[1]["contact_email"], datasetlist[1]["license"], datasetlist[1]["data_standard"], datasetlist[1]["sample_url"], datasetlist[1]["related_articles"], datasetlist[1]["time_saved"], datasetlist[1]["price"], datasetlist[1]["currency_id"], datasetlist[1]["featured_medium_id"]))

		currency.CurrencySelectMock(mock)
		medium.MediumSelectMock(mock)

		tagAssociationSelectMock(mock, 2)

		userDatasetFormatSelectMock(mock, 2)

		delete(datasetlist[1], "tag_ids")

		e.GET(basePath).
			WithQueryObject(map[string]interface{}{
				"limit": "1",
				"page":  "2",
			}).
			WithHeaders(headers).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(datasetlist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(datasetlist[1])

		test.ExpectationsMet(t, mock)
	})
}
