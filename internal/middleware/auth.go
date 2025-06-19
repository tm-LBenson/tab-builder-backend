package middleware

import (
	"context"
	"net/http"
	"strings"

	fbauth "github.com/tm-LBenson/tab-builder-backend/internal/auth"
)

type key int

const uidKey key = 0

func FirebaseAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		uid, err := fbauth.Verify(r.Context(), bearer)
		if err != nil {
			http.Error(w, "unauthenticated: "+err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), uidKey, uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UID(r *http.Request) string {
	return r.Context().Value(uidKey).(string)
}
