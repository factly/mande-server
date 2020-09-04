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
	"gopkg.in/h2non/gock.v1"
)

func TestUpdateMembership(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	test.MeiliGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	e := httpexpect.New(t, server.URL)

	t.Run("update membership", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(MembershipCols).
				AddRow(1, time.Now(), time.Now(), nil, "status", 2, 2, 2))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE \"dp_membership\" SET (.+)  WHERE (.+) \"dp_membership\".\"id\" = `).
			WithArgs(Membership["payment_id"], Membership["plan_id"], Membership["status"], test.AnyTime{}, Membership["user_id"], 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		MembershipSelectMock(mock)

		user.UserSelectMock(mock)

		plan.PlanSelectMock(mock)

		payment.PaymentSelectMock(mock)

		currency.CurrencySelectMock(mock)

		mock.ExpectCommit()

		result := e.PUT(path).
			WithPath("membership_id", "1").
			WithJSON(Membership).
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

		e.PUT(path).
			WithPath("membership_id", "1").
			WithJSON(Membership).
			Expect().
			Status(http.StatusNotFound)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable membership body", func(t *testing.T) {
		e.PUT(path).
			WithPath("membership_id", "1").
			WithJSON(invalidMembership).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("invalid membership id", func(t *testing.T) {
		e.PUT(path).
			WithPath("membership_id", "abc").
			WithJSON(Membership).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("new user does not exist", func(t *testing.T) {
		updateWithErrorExpect(mock, errMembershipUserFK)

		e.PUT(path).
			WithPath("membership_id", "1").
			WithJSON(Membership).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("new plan does not exist", func(t *testing.T) {
		updateWithErrorExpect(mock, errMembershipPlanFK)

		e.PUT(path).
			WithPath("membership_id", "1").
			WithJSON(Membership).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("new payment does not exist", func(t *testing.T) {
		updateWithErrorExpect(mock, errMembershipPaymentFK)

		e.PUT(path).
			WithPath("membership_id", "1").
			WithJSON(Membership).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("update membership when meili is down", func(t *testing.T) {
		gock.Off()
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(MembershipCols).
				AddRow(1, time.Now(), time.Now(), nil, "status", 2, 2, 2))

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE \"dp_membership\" SET (.+)  WHERE (.+) \"dp_membership\".\"id\" = `).
			WithArgs(Membership["payment_id"], Membership["plan_id"], Membership["status"], test.AnyTime{}, Membership["user_id"], 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		MembershipSelectMock(mock)

		user.UserSelectMock(mock)

		plan.PlanSelectMock(mock)

		payment.PaymentSelectMock(mock)

		currency.CurrencySelectMock(mock)

		mock.ExpectRollback()

		e.PUT(path).
			WithPath("membership_id", "1").
			WithJSON(Membership).
			Expect().
			Status(http.StatusInternalServerError)
		test.ExpectationsMet(t, mock)
	})

}
