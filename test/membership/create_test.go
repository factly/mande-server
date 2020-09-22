package membership

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/plan"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestCreateMembership(t *testing.T) {
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

	t.Run("create a membership", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_membership"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, "created", 1, Membership["plan_id"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "payment_id", "razorpay_order_id"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"payment_id", "razorpay_order_id"}).AddRow(nil, nil))

		MembershipSelectMock(mock)

		plan.PlanSelectMock(mock)

		associatedPlansCatalogSelectMock(mock)

		productCatalogAssociationMock(mock, 1)

		mock.ExpectCommit()

		e.POST(basePath).
			WithHeader("X-User", "1").
			WithJSON(requestBody).
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(Membership).
			Value("plan").
			Object().
			ContainsMap(plan.PlanReceive)

		test.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable membership body", func(t *testing.T) {
		e.POST(basePath).
			WithHeader("X-User", "1").
			WithJSON(invalidMembership).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("invalid user header", func(t *testing.T) {
		e.POST(basePath).
			WithHeader("X-User", "abc").
			WithJSON(requestBody).
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("empty membership body", func(t *testing.T) {
		e.POST(basePath).
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("plan does not exist", func(t *testing.T) {
		insertWithErrorExpect(mock, errMembershipPlanFK)

		e.POST(basePath).
			WithHeader("X-User", "1").
			WithJSON(requestBody).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("create a membership when meili is down", func(t *testing.T) {
		gock.Off()
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "dp_membership"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, "created", 1, Membership["plan_id"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "payment_id", "razorpay_order_id"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"payment_id", "razorpay_order_id"}).AddRow(nil, nil))

		MembershipSelectMock(mock)

		plan.PlanSelectMock(mock)

		associatedPlansCatalogSelectMock(mock)

		productCatalogAssociationMock(mock, 1)

		mock.ExpectRollback()

		e.POST(basePath).
			WithHeader("X-User", "1").
			WithJSON(requestBody).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
