package middleware

import (
	"context"
	"net/http"
	"strings"

	fbauth "github.com/tm-LBenson/tab-builder-backend/internal/auth"
	"github.com/tm-LBenson/tab-builder-backend/internal/db"
)

type ctxKey int

const uidKey ctxKey = 0

func FirebaseAuth(store *db.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bearer := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			uid, name, err := fbauth.Verify(r.Context(), bearer)
			if err != nil {
				http.Error(w, "unauthenticated: "+err.Error(), http.StatusUnauthorized)
				return
			}

			if err := store.EnsureUser(r.Context(), uid, name); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), uidKey, uid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UID(r *http.Request) string { return r.Context().Value(uidKey).(string) }
