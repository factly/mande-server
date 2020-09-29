package membership

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/catalog"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/payment"
	"github.com/factly/data-portal-server/test/plan"
	"github.com/gavv/httpexpect"
)

func TestDetailMembership(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	adminExpect := httpexpect.New(t, server.URL)

	// ADMIN tests
	CommonDetailTests(t, mock, adminExpect)

	t.Run("get membership by id", func(t *testing.T) {
		MembershipSelectMock(mock)

		plan.PlanSelectMock(mock)

		catalog.CatalogSelectMock(mock)

		payment.PaymentSelectMock(mock)

		currency.CurrencySelectMock(mock)

		result := adminExpect.GET(path).
			WithPath("membership_id", "1").
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(Membership)

		validateAssociations(result)

		test.ExpectationsMet(t, mock)
	})

	t.Run("membership record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(MembershipCols))

		adminExpect.GET(path).
			WithPath("membership_id", "1").
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	server.Close()

	router = action.RegisterUserRoutes()
	server = httptest.NewServer(router)
	userExpect := httpexpect.New(t, server.URL)

	// USER tests
	CommonDetailTests(t, mock, userExpect)

	t.Run("get membership by id", func(t *testing.T) {
		selectWithTwoArgsMock(mock)

		plan.PlanSelectMock(mock)

		catalog.CatalogSelectMock(mock)

		payment.PaymentSelectMock(mock)

		currency.CurrencySelectMock(mock)

		result := userExpect.GET(path).
			WithPath("membership_id", "1").
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(Membership)

		validateAssociations(result)

		test.ExpectationsMet(t, mock)
	})

	t.Run("membership record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows(MembershipCols))

		userExpect.GET(path).
			WithPath("membership_id", "1").
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid user header", func(t *testing.T) {
		userExpect.GET(path).
			WithPath("membership_id", "1").
			WithHeader("X-User", "abc").
			Expect().
			Status(http.StatusNotFound)
	})

	server.Close()
}

func CommonDetailTests(t *testing.T, mock sqlmock.Sqlmock, e *httpexpect.Expect) {

	t.Run("invalid membership id", func(t *testing.T) {
		e.GET(path).
			WithPath("membership_id", "abc").
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusNotFound)
	})
}
