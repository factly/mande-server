package plan

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/dataset"
	"github.com/factly/data-portal-server/test/tag"
	"github.com/gavv/httpexpect"
)

func TestDetailPlan(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("get plan by id", func(t *testing.T) {
		PlanSelectMock(mock)

		associatedCatalogSelectMock(mock)

		productCatalogAssociationMock(mock, 1)

		currency.CurrencySelectMock(mock)

		dataset.DatasetSelectMock(mock)

		tag.TagSelectMock(mock)

		e.GET(path).
			WithPath("plan_id", "1").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(PlanReceive)

		test.ExpectationsMet(t, mock)
	})

	t.Run("plan record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(PlanCols))

		e.GET(path).
			WithPath("plan_id", "1").
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid id", func(t *testing.T) {
		e.GET(path).
			WithPath("plan_id", "abc").
			Expect().
			Status(http.StatusNotFound)
	})
}
