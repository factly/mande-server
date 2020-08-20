package medium

import "regexp"

var medium map[string]interface{} = map[string]interface{}{
	"name":        "Test Medium",
	"slug":        "test-medium",
	"type":        "testtype",
	"title":       "Test Title",
	"description": "Test Description",
	"caption":     "Test Caption",
	"alt_text":    "Test alt text",
	"file_size":   100,
	"url":         "http:/testurl.com",
	"dimensions":  "testdims",
}

var invalidMedium map[string]interface{} = map[string]interface{}{
	"nam":         "Test Medium",
	"slug":        "test-medium",
	"type":        "testtype",
	"title":       "Test Title",
	"description": "Test Description",
	"caption":     "Test Caption",
	"alt_text":    "Test alt text",
	"filesize":    100,
	"url":         "http:/testurl.com",
	"dimensions":  "testdims",
}

var mediumlist []map[string]interface{} = []map[string]interface{}{
	{
		"name":        "Test Medium 1",
		"slug":        "test-medium-1",
		"type":        "testtype1",
		"title":       "Test Title 1",
		"description": "Test Description 1",
		"caption":     "Test Caption 1",
		"alt_text":    "Test alt text 1",
		"file_size":   100,
		"url":         "http:/testurl1.com",
		"dimensions":  "testdims1",
	},
	{
		"name":        "Test Medium 2",
		"slug":        "test-medium-2",
		"type":        "testtype2",
		"title":       "Test Title 2",
		"description": "Test Description 2",
		"caption":     "Test Caption 2",
		"alt_text":    "Test alt text 2",
		"file_size":   200,
		"url":         "http:/testurl2.com",
		"dimensions":  "testdims2",
	},
}

var mediumCols []string = []string{"id", "created_at", "updated_at", "deleted_at", "name", "slug", "type", "title", "description", "caption", "alt_text", "file_size", "url", "dimensions"}

var selectQuery string = regexp.QuoteMeta(`SELECT * FROM "dp_medium"`)
var countQuery string = regexp.QuoteMeta(`SELECT count(*) FROM "dp_medium"`)
var mediumCatalogQuery string = regexp.QuoteMeta(`SELECT count(*) FROM "dp_catalog`)
var mediumDatasetQuery string = regexp.QuoteMeta(`SELECT count(*) FROM "dp_dataset`)
var mediumProductQuery string = regexp.QuoteMeta(`SELECT count(*) FROM "dp_product`)

var path string = "/media"
var pathId string = "/media/1"
var pathInvalidId string = "/media/abc"
