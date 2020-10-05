package dataset

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/medium"
	"github.com/factly/data-portal-server/test/tag"
	"github.com/gavv/httpexpect"
)

func TestListDataset(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	adminExpect := httpexpect.New(t, server.URL)

	AdminListTests(t, mock, adminExpect)

	server.Close()

	router = action.RegisterUserRoutes()
	server = httptest.NewServer(router)
	userExpect := httpexpect.New(t, server.URL)

	UserListTests(t, mock, userExpect)

	server.Close()
}

func AdminListTests(t *testing.T, mock sqlmock.Sqlmock, e *httpexpect.Expect) {
	t.Run("get empty dataset list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(DatasetCols))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tag" INNER JOIN "dp_dataset_tag"`)).
			WillReturnRows(sqlmock.NewRows(append(tag.TagCols, []string{"tag_id", "dataset_id"}...)))

		e.GET(basePath).
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
				AddRow(1, time.Now(), time.Now(), nil, datasetlist[0]["title"], datasetlist[0]["description"], datasetlist[0]["source"], datasetlist[0]["frequency"], datasetlist[0]["temporal_coverage"], datasetlist[0]["granularity"], datasetlist[0]["contact_name"], datasetlist[0]["contact_email"], datasetlist[0]["license"], datasetlist[0]["data_standard"], datasetlist[0]["sample_url"], datasetlist[0]["related_articles"], datasetlist[0]["time_saved"], datasetlist[0]["price"], datasetlist[0]["currency_id"], datasetlist[0]["featured_medium_id"]).
				AddRow(2, time.Now(), time.Now(), nil, datasetlist[1]["title"], datasetlist[1]["description"], datasetlist[1]["source"], datasetlist[1]["frequency"], datasetlist[1]["temporal_coverage"], datasetlist[1]["granularity"], datasetlist[1]["contact_name"], datasetlist[1]["contact_email"], datasetlist[1]["license"], datasetlist[1]["data_standard"], datasetlist[1]["sample_url"], datasetlist[1]["related_articles"], datasetlist[1]["time_saved"], datasetlist[1]["price"], datasetlist[1]["currency_id"], datasetlist[1]["featured_medium_id"]))

		medium.MediumSelectMock(mock)

		currency.CurrencySelectMock(mock)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tag" INNER JOIN "dp_dataset_tag"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows(append(tag.TagCols, []string{"tag_id", "dataset_id"}...)).
				AddRow(1, time.Now(), time.Now(), nil, "title1", "slug1", 1, 1))

		datasetFormatSelectMock(mock, 1)

		datasetFormatSelectMock(mock, 2)

		delete(datasetlist[0], "tag_ids")

		e.GET(basePath).
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
				AddRow(2, time.Now(), time.Now(), nil, datasetlist[1]["title"], datasetlist[1]["description"], datasetlist[1]["source"], datasetlist[1]["frequency"], datasetlist[1]["temporal_coverage"], datasetlist[1]["granularity"], datasetlist[1]["contact_name"], datasetlist[1]["contact_email"], datasetlist[1]["license"], datasetlist[1]["data_standard"], datasetlist[1]["sample_url"], datasetlist[1]["related_articles"], datasetlist[1]["time_saved"], datasetlist[1]["price"], datasetlist[1]["currency_id"], datasetlist[1]["featured_medium_id"]))

		medium.MediumSelectMock(mock)

		currency.CurrencySelectMock(mock)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tag" INNER JOIN "dp_dataset_tag"`)).
			WithArgs(2).
			WillReturnRows(sqlmock.NewRows(append(tag.TagCols, []string{"tag_id", "dataset_id"}...)).
				AddRow(1, time.Now(), time.Now(), nil, "title1", "slug1", 1, 1))

		datasetFormatSelectMock(mock, 2)

		delete(datasetlist[1], "tag_ids")

		e.GET(basePath).
			WithQueryObject(map[string]interface{}{
				"limit": "1",
				"page":  "2",
			}).
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

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tag" INNER JOIN "dp_dataset_tag"`)).
			WillReturnRows(sqlmock.NewRows(append(tag.TagCols, []string{"tag_id", "dataset_id"}...)))

		e.GET(basePath).
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
				AddRow(1, time.Now(), time.Now(), nil, datasetlist[0]["title"], datasetlist[0]["description"], datasetlist[0]["source"], datasetlist[0]["frequency"], datasetlist[0]["temporal_coverage"], datasetlist[0]["granularity"], datasetlist[0]["contact_name"], datasetlist[0]["contact_email"], datasetlist[0]["license"], datasetlist[0]["data_standard"], datasetlist[0]["sample_url"], datasetlist[0]["related_articles"], datasetlist[0]["time_saved"], datasetlist[0]["price"], datasetlist[0]["currency_id"], datasetlist[0]["featured_medium_id"]).
				AddRow(2, time.Now(), time.Now(), nil, datasetlist[1]["title"], datasetlist[1]["description"], datasetlist[1]["source"], datasetlist[1]["frequency"], datasetlist[1]["temporal_coverage"], datasetlist[1]["granularity"], datasetlist[1]["contact_name"], datasetlist[1]["contact_email"], datasetlist[1]["license"], datasetlist[1]["data_standard"], datasetlist[1]["sample_url"], datasetlist[1]["related_articles"], datasetlist[1]["time_saved"], datasetlist[1]["price"], datasetlist[1]["currency_id"], datasetlist[1]["featured_medium_id"]))

		medium.MediumSelectMock(mock)

		currency.CurrencySelectMock(mock)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tag" INNER JOIN "dp_dataset_tag"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows(append(tag.TagCols, []string{"tag_id", "dataset_id"}...)).
				AddRow(1, time.Now(), time.Now(), nil, "title1", "slug1", 1, 1))

		userDatasetFormatSelectMock(mock, 1)

		userDatasetFormatSelectMock(mock, 2)

		delete(datasetlist[0], "tag_ids")

		e.GET(basePath).
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
				AddRow(2, time.Now(), time.Now(), nil, datasetlist[1]["title"], datasetlist[1]["description"], datasetlist[1]["source"], datasetlist[1]["frequency"], datasetlist[1]["temporal_coverage"], datasetlist[1]["granularity"], datasetlist[1]["contact_name"], datasetlist[1]["contact_email"], datasetlist[1]["license"], datasetlist[1]["data_standard"], datasetlist[1]["sample_url"], datasetlist[1]["related_articles"], datasetlist[1]["time_saved"], datasetlist[1]["price"], datasetlist[1]["currency_id"], datasetlist[1]["featured_medium_id"]))

		medium.MediumSelectMock(mock)

		currency.CurrencySelectMock(mock)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_tag" INNER JOIN "dp_dataset_tag"`)).
			WithArgs(2).
			WillReturnRows(sqlmock.NewRows(append(tag.TagCols, []string{"tag_id", "dataset_id"}...)).
				AddRow(1, time.Now(), time.Now(), nil, "title1", "slug1", 1, 1))

		userDatasetFormatSelectMock(mock, 2)

		delete(datasetlist[1], "tag_ids")

		e.GET(basePath).
			WithQueryObject(map[string]interface{}{
				"limit": "1",
				"page":  "2",
			}).
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
