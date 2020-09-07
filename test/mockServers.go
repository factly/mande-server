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
			"object_id":   "format_2",
			"kind":        "format",
			"id":          2,
			"description": "This is a test format",
			"name":        "Test Format",
			"is_default":  true,
		},
		{
			"object_id":   "format_3",
			"kind":        "format",
			"id":          3,
			"description": "This is second test format",
			"name":        "Test format 2",
			"is_default":  true,
		},
		{
			"object_id": "tag_2",
			"kind":      "tag",
			"id":        2,
			"slug":      "test-tag",
			"title":     "Test tag",
		},
	},
	"offset":           0,
	"limit":            10,
	"nbHits":           10,
	"exhaustiveNbHits": false,
	"processingTimeMs": 2,
	"query":            "test",
}

func MeiliGock() {
	gock.New(config.MeiliURL).
		Post("/indexes/data-portal/search").
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
