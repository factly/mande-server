package membership

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/catalog"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/payment"
	"github.com/factly/data-portal-server/test/plan"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestListMembership(t *testing.T) {

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
	CommonListTests(t, mock, adminExpect)

	t.Run("get memberships list with user query parameter", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(membershiplist)))

		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(MembershipCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, membershiplist[0]["status"], membershiplist[0]["user_id"], membershiplist[0]["payment_id"], membershiplist[0]["plan_id"], membershiplist[0]["razorpay_order_id"]).
				AddRow(2, time.Now(), time.Now(), nil, 1, 1, membershiplist[1]["status"], membershiplist[1]["user_id"], membershiplist[1]["payment_id"], membershiplist[1]["plan_id"], membershiplist[1]["razorpay_order_id"]))

		payment.PaymentSelectMock(mock)
		currency.CurrencySelectMock(mock)

		plan.PlanSelectMock(mock)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_plan_catalog"`)).
			WillReturnRows(sqlmock.NewRows([]string{"plan_id", "catalog_id"}).
				AddRow(1, 1))

		catalog.CatalogSelectMock(mock)

		adminExpect.GET(basePath).
			WithHeaders(headers).
			WithQuery("user", "1").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(membershiplist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(membershiplist[0])

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid user query param", func(t *testing.T) {
		adminExpect.GET(basePath).
			WithHeaders(headers).
			WithQuery("user", "abc").
			Expect().
			Status(http.StatusBadRequest)
	})

	server.Close()

	router = action.RegisterUserRoutes()
	server = httptest.NewServer(router)
	userExpect := httpexpect.New(t, server.URL)

	gock.New(server.URL).EnableNetworking().Persist()

	// USER tests
	CommonListTests(t, mock, userExpect)

	server.Close()
}

func CommonListTests(t *testing.T, mock sqlmock.Sqlmock, e *httpexpect.Expect) {
	t.Run("get empty memberships list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(MembershipCols))

		e.GET(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 0})

		test.ExpectationsMet(t, mock)
	})

	t.Run("get memberships list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(len(membershiplist)))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(MembershipCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, membershiplist[0]["status"], membershiplist[0]["user_id"], membershiplist[0]["payment_id"], membershiplist[0]["plan_id"], membershiplist[0]["razorpay_order_id"]).
				AddRow(2, time.Now(), time.Now(), nil, 1, 1, membershiplist[1]["status"], membershiplist[1]["user_id"], membershiplist[1]["payment_id"], membershiplist[1]["plan_id"], membershiplist[1]["razorpay_order_id"]))

		payment.PaymentSelectMock(mock)
		currency.CurrencySelectMock(mock)

		plan.PlanSelectMock(mock)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_plan_catalog"`)).
			WillReturnRows(sqlmock.NewRows([]string{"plan_id", "catalog_id"}).
				AddRow(1, 1))

		catalog.CatalogSelectMock(mock)

		e.GET(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": len(membershiplist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(membershiplist[0])

		test.ExpectationsMet(t, mock)
	})

	t.Run("get memberships list with paiganation", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("2"))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(MembershipCols).
				AddRow(2, time.Now(), time.Now(), nil, 1, 1, membershiplist[1]["status"], membershiplist[1]["user_id"], membershiplist[1]["payment_id"], membershiplist[1]["plan_id"], membershiplist[1]["razorpay_order_id"]))

		payment.PaymentSelectMock(mock)
		currency.CurrencySelectMock(mock)

		plan.PlanSelectMock(mock)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_plan_catalog"`)).
			WillReturnRows(sqlmock.NewRows([]string{"plan_id", "catalog_id"}).
				AddRow(1, 1))

		catalog.CatalogSelectMock(mock)

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
			ContainsMap(map[string]interface{}{"total": len(membershiplist)}).
			Value("nodes").
			Array().
			Element(0).
			Object().
			ContainsMap(membershiplist[1])

		test.ExpectationsMet(t, mock)
	})
}
