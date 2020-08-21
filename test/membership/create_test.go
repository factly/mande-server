package membership

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/payment"
	"github.com/factly/data-portal-server/test/plan"
	"github.com/factly/data-portal-server/test/user"
	"github.com/gavv/httpexpect"
)

func TestCreateMembership(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	t.Run("create a membership", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_membership"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, Membership["status"], Membership["user_id"], Membership["payment_id"], Membership["plan_id"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		MembershipSelectMock(mock)

		user.UserSelectMock(mock)

		plan.PlanSelectMock(mock)

		payment.PaymentSelectMock(mock)

		currency.CurrencySelectMock(mock)

		result := e.POST(basePath).
			WithJSON(Membership).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(Membership)

		validateAssociations(result)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable membership body", func(t *testing.T) {
		e.POST(basePath).
			WithJSON(invalidMembership).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("empty membership body", func(t *testing.T) {
		e.POST(basePath).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("user does not exist", func(t *testing.T) {
		insertWithErrorExpect(mock, errMembershipUserFK)

		e.POST(basePath).
			WithJSON(Membership).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("plan does not exist", func(t *testing.T) {
		insertWithErrorExpect(mock, errMembershipPlanFK)

		e.POST(basePath).
			WithJSON(Membership).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("payment does not exist", func(t *testing.T) {
		insertWithErrorExpect(mock, errMembershipPaymentFK)

		e.POST(basePath).
			WithJSON(Membership).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})
}
