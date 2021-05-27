package currency

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/mande-server/action"
	"github.com/factly/mande-server/action/currency"
	"github.com/factly/mande-server/test"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestCreateDefaultCurrency(t *testing.T) {
	// Setup DB
	mock := test.SetupMockDB()

	currency.DataFile = "../../data/currency.json"

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	test.KetoGock()
	test.KavachGock()
	gock.New(server.URL).EnableNetworking().Persist()

	t.Run("create default currency", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(CurrencyCols))

		mock.ExpectQuery(`INSERT INTO "dp_currency"`).
			WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, Currency["iso_code"], Currency["name"]).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
		mock.ExpectCommit()

		e.POST(defaultPath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusCreated).JSON().
			Object().
			Value("nodes").
			Array().Element(0).Object().ContainsMap(Currency)
		test.ExpectationsMet(t, mock)
	})

	t.Run("default currency already created", func(t *testing.T) {
		mock.ExpectBegin()
		CurrencySelectMock(mock)

		mock.ExpectCommit()

		e.POST(defaultPath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusCreated).JSON().
			Object().
			Value("nodes").
			Array().Element(0).Object().ContainsMap(Currency)
		test.ExpectationsMet(t, mock)
	})

	t.Run("cannot open data file", func(t *testing.T) {
		currency.DataFile = "nofile.json"
		e.POST(defaultPath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusInternalServerError)
		currency.DataFile = "../../data/currency.json"
	})

	t.Run("cannot parse data file", func(t *testing.T) {
		currency.DataFile = "invalidData.json"
		e.POST(defaultPath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusInternalServerError)
		currency.DataFile = "../../data/currency.json"
	})
}
