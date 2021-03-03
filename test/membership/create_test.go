package membership

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/mande-server/action"
	"github.com/factly/mande-server/test"
	"github.com/factly/mande-server/test/currency"
	"github.com/factly/mande-server/test/plan"
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
	test.KetoGock()
	test.KavachGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	e := httpexpect.New(t, server.URL)

	t.Run("create a membership", func(t *testing.T) {
		mock.ExpectBegin()
		createMock(mock)
		mock.ExpectCommit()

		e.POST(basePath).
			WithHeaders(headers).
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
			WithHeaders(headers).
			WithJSON(invalidMembership).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("empty membership body", func(t *testing.T) {
		e.POST(basePath).
			WithHeaders(headers).
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
			WithHeaders(headers).
			WithJSON(requestBody).
			Expect().
			Status(http.StatusUnprocessableEntity)

		test.ExpectationsMet(t, mock)
	})

	t.Run("creating membership fails", func(t *testing.T) {
		mock.ExpectBegin()
		plan.PlanSelectMock(mock)
		mock.ExpectQuery(`INSERT INTO "dp_membership"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, "created", 1, nil, Membership["plan_id"], "").
			WillReturnError(errors.New(`cannot create membership`))

		mock.ExpectRollback()
		e.POST(basePath).
			WithHeaders(headers).
			WithJSON(requestBody).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("razorpay gives error response", func(t *testing.T) {
		gock.Off()
		test.MeiliGock()
		test.KavachGock()
		test.KetoGock()
		gock.New(server.URL).EnableNetworking().Persist()
		defer gock.DisableNetworking()

		gock.New("https://api.razorpay.com").
			Post("/v1/orders").
			Persist().
			Reply(http.StatusInternalServerError)

		mock.ExpectBegin()
		plan.PlanSelectMock(mock)
		mock.ExpectQuery(`INSERT INTO "dp_membership"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, "created", 1, nil, Membership["plan_id"], "").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_currency"`)).
			WillReturnRows(sqlmock.NewRows(currency.CurrencyCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, currency.Currency["iso_code"], currency.Currency["name"]))

		mock.ExpectRollback()
		e.POST(basePath).
			WithHeaders(headers).
			WithJSON(requestBody).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("razorpay returns invalid body", func(t *testing.T) {
		gock.Off()
		test.MeiliGock()
		test.KavachGock()
		test.KetoGock()
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
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, "created", 1, nil, Membership["plan_id"], "").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_currency"`)).
			WillReturnRows(sqlmock.NewRows(currency.CurrencyCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, currency.Currency["iso_code"], currency.Currency["name"]))

		mock.ExpectRollback()
		e.POST(basePath).
			WithHeaders(headers).
			WithJSON(requestBody).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	gock.Off()
	test.MeiliGock()
	test.RazorpayGock()
	test.KavachGock()
	test.KetoGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	t.Run("updating membership fails", func(t *testing.T) {
		mock.ExpectBegin()
		plan.PlanSelectMock(mock)
		mock.ExpectQuery(`INSERT INTO "dp_membership"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, "created", 1, nil, Membership["plan_id"], "").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "dp_currency"`)).
			WillReturnRows(sqlmock.NewRows(currency.CurrencyCols).
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, currency.Currency["iso_code"], currency.Currency["name"]))

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "dp_membership" SET`)).
			WithArgs(1, test.AnyTime{}, test.AnyTime{}, 1, 1, "processing", 1, 1, test.RazorpayOrder["id"], 1).
			WillReturnError(errors.New(`cannot update membership`))

		mock.ExpectRollback()
		e.POST(basePath).
			WithHeaders(headers).
			WithJSON(requestBody).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

	t.Run("create a membership when meili is down", func(t *testing.T) {
		gock.Off()
		test.RazorpayGock()
		test.KavachGock()
		test.KetoGock()
		gock.New(server.URL).EnableNetworking().Persist()

		mock.ExpectBegin()
		createMock(mock)
		mock.ExpectRollback()

		e.POST(basePath).
			WithHeaders(headers).
			WithJSON(requestBody).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
