package order

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestCreateOrder(t *testing.T) {
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

	t.Run("create a order", func(t *testing.T) {

		insertMock(mock)
		mock.ExpectCommit()

		result := e.POST(basePath).
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusCreated).
			JSON().
			Object().
			ContainsMap(Order)

		validateAssociations(result)

		test.ExpectationsMet(t, mock)
	})

	t.Run("invalid user header", func(t *testing.T) {
		e.POST(basePath).
			WithHeader("X-User", "abc").
			Expect().
			Status(http.StatusNotFound)
	})

	t.Run("create a order when meili is down", func(t *testing.T) {
		gock.Off()
		insertMock(mock)
		mock.ExpectRollback()

		e.POST(basePath).
			WithHeader("X-User", "1").
			WithJSON(Order).
			Expect().
			Status(http.StatusInternalServerError)

		test.ExpectationsMet(t, mock)
	})

}
