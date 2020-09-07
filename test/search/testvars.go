package search

var Data = map[string]interface{}{
	"q":       "test",
	"limit":   10,
	"filters": "kind=tag OR kind=format",
}

var undecodableData = map[string]interface{}{
	"q":     10,
	"limit": "10",
}

var invalidData = map[string]interface{}{
	"q":       "te",
	"limit":   100,
	"filters": "kind=category",
}

var path string = "/search"
