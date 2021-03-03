package search

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/factly/mande-server/action"
	"github.com/factly/mande-server/test"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestListSearch(t *testing.T) {
	// Setup DB
	test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterAdminRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	test.MeiliGock()
	test.KavachGock()
	test.KetoGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	e := httpexpect.New(t, server.URL)

	t.Run("undecodable body", func(t *testing.T) {
		e.POST(path).
			WithJSON(undecodableData).
			WithHeaders(headers).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("invalid body", func(t *testing.T) {
		e.POST(path).
			WithJSON(invalidData).
			WithHeaders(headers).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("search entities with query 'test'", func(t *testing.T) {
		e.POST(path).
			WithJSON(Data).
			WithHeaders(headers).
			Expect().
			Status(http.StatusOK)
	})

	t.Run("meili server is down", func(t *testing.T) {
		gock.Off()
		test.KavachGock()
		test.KetoGock()
		gock.New(server.URL).EnableNetworking().Persist()

		e.POST(path).
			WithJSON(Data).
			WithHeaders(headers).
			Expect().
			Status(http.StatusInternalServerError)
	})
}
