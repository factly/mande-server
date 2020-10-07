package test

import (
	"net/http"

	"github.com/factly/data-portal-server/util/razorpay"

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

var RazorpayOrder = map[string]interface{}{
	"id":          "order_FltCdu23fGaTwG",
	"entity":      "order",
	"amount":      5000,
	"amount_paid": 0,
	"amount_due":  5000,
	"currency":    "INR",
	"receipt":     "Test Receipt no. 1",
	"offer_id":    nil,
	"status":      "created",
	"attempts":    0,
	"notes": map[string]interface{}{
		"info": "this payment is for first order",
	},
	"created_at": 1602047090,
}

func RazorpayGock() {
	razorpay.SetupClient()

	gock.New("https://api.razorpay.com").
		Post("/v1/orders").
		Persist().
		Reply(http.StatusOK).
		JSON(RazorpayOrder)
}
