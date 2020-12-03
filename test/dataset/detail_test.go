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
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestDetailDataset(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	adminExpect := httpexpect.New(t, server.URL)

	test.KetoGock()
	test.KavachGock()
	gock.New(server.URL).EnableNetworking().Persist()

	// ADMIN tests
	AdminTests(t, mock, adminExpect)

	server.Close()

	router = action.RegisterUserRoutes()
	server = httptest.NewServer(router)
	userExpect := httpexpect.New(t, server.URL)

	gock.New(server.URL).EnableNetworking().Persist()
	// USER tests
	UserTests(t, mock, userExpect)

	server.Close()
}

func AdminTests(t *testing.T, mock sqlmock.Sqlmock, e *httpexpect.Expect) {
	t.Run("get dataset by id", func(t *testing.T) {
		DatasetSelectMock(mock)

		currency.CurrencySelectMock(mock)
		medium.MediumSelectMock(mock)

		tagAssociationSelectMock(mock)

		datasetFormatSelectMock(mock, 1)

		result := e.GET(path).
			WithPath("dataset_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(DatasetReceive)

		validateAssociations(result)

		test.ExpectationsMet(t, mock)
	})

	t.Run("dataset record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(DatasetCols))

		e.GET(path).
			WithHeaders(headers).
			WithPath("dataset_id", "1").
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid dataset id", func(t *testing.T) {
		e.GET(path).
			WithPath("dataset_id", "abc").
			WithHeaders(headers).
			Expect().
			Status(http.StatusBadRequest)
	})
}

func UserTests(t *testing.T, mock sqlmock.Sqlmock, e *httpexpect.Expect) {
	t.Run("get dataset by id", func(t *testing.T) {
		DatasetSelectMock(mock)

		currency.CurrencySelectMock(mock)
		medium.MediumSelectMock(mock)

		tagAssociationSelectMock(mock)

		orderSelectMock(mock)

		mock.ExpectQuery(`SELECT "dp_product"(.+) INNER JOIN dp_order_item`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(productCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, "title", "slug", 100, "status", 1, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "dp_product" JOIN "dp_product_dataset"`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		datasetFormatSelectMock(mock, 1)

		result := e.GET(path).
			WithPath("dataset_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(DatasetReceive)

		validateAssociations(result)

		test.ExpectationsMet(t, mock)
	})

	t.Run("get dataset when user does not own dataset", func(t *testing.T) {
		DatasetSelectMock(mock)

		currency.CurrencySelectMock(mock)
		medium.MediumSelectMock(mock)

		tagAssociationSelectMock(mock)

		orderSelectMock(mock)

		mock.ExpectQuery(`SELECT "dp_product"(.+) INNER JOIN dp_order_item`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(productCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, "title", "slug", 100, "status", 1, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "dp_product" JOIN "dp_product_dataset"`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		result := e.GET(path).
			WithPath("dataset_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(DatasetReceive)

		result.Value("formats").
			Array().
			Length().
			Equal(0)

		validateAssociations(result)

		test.ExpectationsMet(t, mock)
	})

	t.Run("dataset record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(DatasetCols))

		e.GET(path).
			WithPath("dataset_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid dataset id", func(t *testing.T) {
		e.GET(path).
			WithPath("dataset_id", "abc").
			WithHeaders(headers).
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("invalid user id header", func(t *testing.T) {
		e.GET(path).
			WithPath("dataset_id", "1").
			WithHeaders(headers).
			Expect().
			Status(http.StatusNotFound)
	})
}
