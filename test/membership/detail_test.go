package membership

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/mande-server/action"
	"github.com/factly/mande-server/test"
	"github.com/factly/mande-server/test/catalog"
	"github.com/factly/mande-server/test/currency"
	"github.com/factly/mande-server/test/payment"
	"github.com/factly/mande-server/test/plan"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestDetailMembership(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	adminExpect := httpexpect.New(t, server.URL)

	test.KetoGock()
	test.KavachGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	// ADMIN tests
	CommonDetailTests(t, mock, adminExpect)

	t.Run("get membership by id", func(t *testing.T) {
		MembershipSelectMock(mock)

		payment.PaymentSelectMock(mock)
		currency.CurrencySelectMock(mock)

		plan.PlanSelectMock(mock)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_plan_catalog"`)).
			WillReturnRows(sqlmock.NewRows([]string{"plan_id", "catalog_id"}).
				AddRow(1, 1))

		catalog.CatalogSelectMock(mock)

		result := adminExpect.GET(path).
			WithPath("membership_id", "1").
			WithHeaders(headers).
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
			WithHeaders(headers).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	server.Close()

	router = action.RegisterUserRoutes()
	server = httptest.NewServer(router)
	userExpect := httpexpect.New(t, server.URL)

	gock.New(server.URL).EnableNetworking().Persist()

	// USER tests
	CommonDetailTests(t, mock, userExpect)

	t.Run("get membership by id", func(t *testing.T) {
		MembershipSelectMock(mock, 1, 1)

		payment.PaymentSelectMock(mock)
		currency.CurrencySelectMock(mock)

		plan.PlanSelectMock(mock)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_plan_catalog"`)).
			WillReturnRows(sqlmock.NewRows([]string{"plan_id", "catalog_id"}).
				AddRow(1, 1))

		catalog.CatalogSelectMock(mock)

		result := userExpect.GET(path).
			WithPath("membership_id", "1").
			WithHeaders(headers).
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
			WithHeaders(headers).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	server.Close()
}

func CommonDetailTests(t *testing.T, mock sqlmock.Sqlmock, e *httpexpect.Expect) {

	t.Run("invalid membership id", func(t *testing.T) {
		e.GET(path).
			WithPath("membership_id", "abc").
			WithHeaders(headers).
			Expect().
			Status(http.StatusBadRequest)
	})
}
