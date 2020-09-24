package dataset

import (
	"net/http"
	"net/http/httptest"
	"testing"

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

	CommonDetailTests(t, mock, adminExpect)

	server.Close()

	router = action.RegisterUserRoutes()
	server = httptest.NewServer(router)
	userExpect := httpexpect.New(t, server.URL)

	CommonDetailTests(t, mock, userExpect)

	server.Close()
}

func CommonDetailTests(t *testing.T, mock sqlmock.Sqlmock, e *httpexpect.Expect) {
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
