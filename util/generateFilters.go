package util

import (
	"fmt"

	"github.com/factly/x/meilisearchx"
)

func GenerateFilters(tagIDs []string) string {
	filters := ""
	if len(tagIDs) > 0 {
		filters = fmt.Sprint(filters, meilisearchx.GenerateFieldFilter(tagIDs, "tag_ids"), " AND ")
	}

	if filters != "" && filters[len(filters)-5:] == " AND " {
		filters = filters[:len(filters)-5]
	}

	return filters
}
