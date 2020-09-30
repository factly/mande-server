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
)

func TestDetailDataset(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	adminExpect := httpexpect.New(t, server.URL)

	// ADMIN tests
	AdminTests(t, mock, adminExpect)

	server.Close()

	router = action.RegisterUserRoutes()
	server = httptest.NewServer(router)
	userExpect := httpexpect.New(t, server.URL)

	// USER tests
	UserTests(t, mock, userExpect)

	server.Close()
}

func AdminTests(t *testing.T, mock sqlmock.Sqlmock, e *httpexpect.Expect) {
	t.Run("get dataset by id", func(t *testing.T) {
		DatasetSelectMock(mock)

		medium.MediumSelectMock(mock)

		currency.CurrencySelectMock(mock)

		tagAssociationSelectMock(mock)

		datasetFormatSelectMock(mock, 1)

		result := e.GET(path).
			WithPath("dataset_id", "1").
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
			WithPath("dataset_id", "1").
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid dataset id", func(t *testing.T) {
		e.GET(path).
			WithPath("dataset_id", "abc").
			Expect().
			Status(http.StatusNotFound)
	})
}

func UserTests(t *testing.T, mock sqlmock.Sqlmock, e *httpexpect.Expect) {
	t.Run("get dataset by id", func(t *testing.T) {
		DatasetSelectMock(mock)

		medium.MediumSelectMock(mock)

		currency.CurrencySelectMock(mock)

		tagAssociationSelectMock(mock)

		orderSelectMock(mock)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "dp_product".* FROM "dp_product" INNER JOIN dp_order_item`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(productCols).
				AddRow(1, time.Now(), time.Now(), nil, "title", "slug", 100, "status", 1, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "dp_product" INNER JOIN "dp_product_dataset"`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		datasetFormatSelectMock(mock, 1)

		result := e.GET(path).
			WithPath("dataset_id", "1").
			WithHeader("X-User", "1").
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

		medium.MediumSelectMock(mock)

		currency.CurrencySelectMock(mock)

		tagAssociationSelectMock(mock)

		orderSelectMock(mock)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "dp_product".* FROM "dp_product" INNER JOIN dp_order_item`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(productCols).
				AddRow(1, time.Now(), time.Now(), nil, "title", "slug", 100, "status", 1, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "dp_product" INNER JOIN "dp_product_dataset"`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		result := e.GET(path).
			WithPath("dataset_id", "1").
			WithHeader("X-User", "1").
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
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid dataset id", func(t *testing.T) {
		e.GET(path).
			WithPath("dataset_id", "abc").
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("invalid user id header", func(t *testing.T) {
		e.GET(path).
			WithPath("dataset_id", "1").
			WithHeader("X-User", "abc").
			Expect().
			Status(http.StatusNotFound)
	})
}
