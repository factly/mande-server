package util

import (
	"context"
	"errors"
	"net/http"
	"strconv"
)

type ctxKeyOrganisationID int

// OrganisationIDKey is the key that holds the unique Organisation ID in a request context.
const OrganisationIDKey ctxKeyOrganisationID = 0

// CheckOrganisation check X-Organisation in header
func CheckOrganisation(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		org := r.Header.Get("X-Organisation")
		if org == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		oid, err := strconv.Atoi(org)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, OrganisationIDKey, oid)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetOrganisation return Organisation ID
func GetOrganisation(ctx context.Context) (int, error) {
	if ctx == nil {
		return 0, errors.New("context not found")
	}
	orgID := ctx.Value(OrganisationIDKey)
	if orgID != nil {
		return orgID.(int), nil
	}
	return 0, errors.New("something went wrong")
}
