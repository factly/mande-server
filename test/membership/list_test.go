package membership

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/payment"
	"github.com/factly/data-portal-server/test/plan"
	"github.com/factly/data-portal-server/test/user"
	"github.com/gavv/httpexpect"
)

func TestListMembership(t *testing.T) {

	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("get empty memberships list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("0"))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(MembershipCols))

		e.GET(basePath).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 0})

		test.ExpectationsMet(t, mock)
	})

	t.Run("get memberships list", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("2"))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(MembershipCols).
				AddRow(1, time.Now(), time.Now(), nil, membershiplist[0]["status"], membershiplist[0]["user_id"], membershiplist[0]["payment_id"], membershiplist[0]["plan_id"]).
				AddRow(2, time.Now(), time.Now(), nil, membershiplist[1]["status"], membershiplist[1]["user_id"], membershiplist[1]["payment_id"], membershiplist[1]["plan_id"]))

		user.UserSelectMock(mock)

		plan.PlanSelectMock(mock)

		payment.PaymentSelectMock(mock)

		currency.CurrencySelectMock(mock)

		e.GET(basePath).
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
				AddRow(2, time.Now(), time.Now(), nil, membershiplist[1]["status"], membershiplist[1]["user_id"], membershiplist[1]["payment_id"], membershiplist[1]["plan_id"]))

		user.UserSelectMock(mock)

		plan.PlanSelectMock(mock)

		payment.PaymentSelectMock(mock)

		currency.CurrencySelectMock(mock)

		e.GET(basePath).
			WithQueryObject(map[string]interface{}{
				"limit": "1",
				"page":  "2",
			}).
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
