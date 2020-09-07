package search

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/factly/data-portal-server/action"
	"github.com/factly/data-portal-server/test"
	"github.com/gavv/httpexpect"
	"gopkg.in/h2non/gock.v1"
)

func TestListSearch(t *testing.T) {
	// Setup DB
	test.SetupMockDB()

	// Setup HttpExpect
	router := action.RegisterRoutes()
	server := httptest.NewServer(router)
	defer server.Close()

	test.MeiliGock()
	gock.New(server.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()

	e := httpexpect.New(t, server.URL)

	t.Run("undecodable body", func(t *testing.T) {
		e.POST(path).
			WithJSON(undecodableData).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("invalid body", func(t *testing.T) {
		e.POST(path).
			WithJSON(invalidData).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("search entities with query 'test'", func(t *testing.T) {
		e.POST(path).
			WithJSON(Data).
			Expect().
			Status(http.StatusOK)
	})

	t.Run("meili server is down", func(t *testing.T) {
		gock.Off()
		e.POST(path).
			WithJSON(Data).
			Expect().
			Status(http.StatusInternalServerError)
	})
}
