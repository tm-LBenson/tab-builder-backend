package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/you/tabs-api/internal/db"
)

type SongHandler struct{ Store *db.Store }

func (h SongHandler) Register(r chi.Router) {
	r.With(h.withUID).Get("/", h.list)
	// more routes later
}

func (h SongHandler) withUID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value("uid") == nil {
			http.Error(w, "unauthenticated", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h SongHandler) list(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value("uid").(string)
	rows, err := h.Store.ListSongs(r.Context(), uid)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(rows)
}
