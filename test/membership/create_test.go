package membership

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/factly/data-portal-server/test/currency"
	"github.com/factly/data-portal-server/test/plan"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestCreateMembership(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterUserRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	test.MeiliGock()
	test.RazorpayGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	e := httpexpect.New(t, server.URL)

	t.Run("create a membership", func(t *testing.T) {
		mock.ExpectBegin()
		createMock(mock)
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
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_plan"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(plan.PlanCols))

		mock.ExpectRollback()
		e.POST(basePath).
			WithHeader("X-User", "1").
			WithJSON(requestBody).
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})

	t.Run("creating membership fails", func(t *testing.T) {
		mock.ExpectBegin()
		plan.PlanSelectMock(mock)
		mock.ExpectQuery(`INSERT INTO "dp_membership"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, "created", 1, Membership["plan_id"]).
			WillReturnError(errors.New(`cannot create membership`))

		mock.ExpectRollback()
		e.POST(basePath).
			WithHeader("X-User", "1").
			WithJSON(requestBody).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("razorpay gives error response", func(t *testing.T) {
		gock.Off()
		test.MeiliGock()
		gock.New(server.URL).EnableNetworking().Persist()
		defer gock.DisableNetworking()

		gock.New("https://api.razorpay.com").
			Post("/v1/orders").
			Persist().
			Reply(http.StatusInternalServerError)

		mock.ExpectBegin()
		plan.PlanSelectMock(mock)
		mock.ExpectQuery(`INSERT INTO "dp_membership"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, "created", 1, Membership["plan_id"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "payment_id", "razorpay_order_id"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"payment_id", "razorpay_order_id"}).AddRow(nil, nil))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_currency"`)).
			WillReturnRows(sqlmock.NewRows(currency.CurrencyCols).
				AddRow(1, time.Now(), time.Now(), nil, currency.Currency["iso_code"], currency.Currency["name"]))

		mock.ExpectRollback()
		e.POST(basePath).
			WithHeader("X-User", "1").
			WithJSON(requestBody).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("razorpay returns invalid body", func(t *testing.T) {
		gock.Off()
		test.MeiliGock()
		gock.New(server.URL).EnableNetworking().Persist()
		defer gock.DisableNetworking()

		gock.New("https://api.razorpay.com").
			Post("/v1/orders").
			Persist().
			Reply(http.StatusOK).
			JSON(map[string]interface{}{
				"amount":      5000,
				"amount_paid": 0,
				"amount_due":  5000,
				"currency":    "INR",
				"receipt":     "Test Receipt no. 1",
			})

		mock.ExpectBegin()
		plan.PlanSelectMock(mock)
		mock.ExpectQuery(`INSERT INTO "dp_membership"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, "created", 1, Membership["plan_id"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "payment_id", "razorpay_order_id"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"payment_id", "razorpay_order_id"}).AddRow(nil, nil))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_currency"`)).
			WillReturnRows(sqlmock.NewRows(currency.CurrencyCols).
				AddRow(1, time.Now(), time.Now(), nil, currency.Currency["iso_code"], currency.Currency["name"]))

		mock.ExpectRollback()
		e.POST(basePath).
			WithHeader("X-User", "1").
			WithJSON(requestBody).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	gock.Off()
	test.MeiliGock()
	test.RazorpayGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	t.Run("updating membership fails", func(t *testing.T) {
		mock.ExpectBegin()
		plan.PlanSelectMock(mock)
		mock.ExpectQuery(`INSERT INTO "dp_membership"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, "created", 1, Membership["plan_id"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT "payment_id", "razorpay_order_id"`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"payment_id", "razorpay_order_id"}).AddRow(nil, nil))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_currency"`)).
			WillReturnRows(sqlmock.NewRows(currency.CurrencyCols).
				AddRow(1, time.Now(), time.Now(), nil, currency.Currency["iso_code"], currency.Currency["name"]))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_membership" SET`)).
			WithArgs(test.AnyTime{}, 1, 1, test.RazorpayOrder["id"], "processing", test.AnyTime{}, 1, 1).
			WillReturnError(errors.New(`cannot update membership`))

		mock.ExpectRollback()
		e.POST(basePath).
			WithHeader("X-User", "1").
			WithJSON(requestBody).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("create a membership when meili is down", func(t *testing.T) {
		gock.Off()
		test.RazorpayGock()
		gock.New(server.URL).EnableNetworking().Persist()

		mock.ExpectBegin()
		createMock(mock)
		mock.ExpectRollback()

		e.POST(basePath).
			WithHeader("X-User", "1").
			WithJSON(requestBody).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
