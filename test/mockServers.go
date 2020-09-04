package test

import (
	"net/http"

	"github.com/factly/data-portal-server/config"
	"gopkg.in/h2non/gock.v1"
)

var ReturnUpdate = map[string]interface{}{
	"updateId": 1,
}

var MeiliHits = map[string]interface{}{
	"hits": []map[string]interface{}{
		{
			"description": "Test claimant",
			"id":          2,
			"kind":        "claimant",
			"medium_id":   2,
			"name":        "Tester",
			"object_id":   "claimant_2",
			"slug":        "tester",
			"space_id":    1,
			"tag_line":    "A claimant for testing",
		},
		{
			"description": "this is test category",
			"id":          3,
			"kind":        "category",
			"medium_id":   2,
			"name":        "Test category",
			"object_id":   "category_3",
			"slug":        "test-category",
			"space_id":    1,
		},
	},
	"offset":           0,
	"limit":            20,
	"nbHits":           7,
	"exhaustiveNbHits": false,
	"processingTimeMs": 2,
	"query":            "test",
}

func MeiliGock() {
	gock.New(config.MeiliURL + "/indexes/data-portal/search").
		HeaderPresent("X-Meili-API-Key").
		Persist().
		Reply(http.StatusOK).
		JSON(MeiliHits)

	gock.New(config.MeiliURL).
		Post("/indexes/data-portal/documents").
		HeaderPresent("X-Meili-API-Key").
		Persist().
		Reply(http.StatusAccepted).
		JSON(ReturnUpdate)

	gock.New(config.MeiliURL).
		Put("/indexes/data-portal/documents").
		HeaderPresent("X-Meili-API-Key").
		Persist().
		Reply(http.StatusAccepted).
		JSON(ReturnUpdate)

	gock.New(config.MeiliURL).
		Delete("/indexes/data-portal/documents/(.+)").
		HeaderPresent("X-Meili-API-Key").
		Persist().
		Reply(http.StatusAccepted).
		JSON(ReturnUpdate)
}
