package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tm-LBenson/tab-builder-backend/internal/db"
	"github.com/tm-LBenson/tab-builder-backend/internal/middleware"
)

type SongHandler struct{ Store *db.Store }

func (h SongHandler) Register(r chi.Router) {
	r.Get("/", h.list)
}

func (h SongHandler) list(w http.ResponseWriter, r *http.Request) {
	uid := middleware.UID(r)
	rows, err := h.Store.ListSongs(r.Context(), uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(rows)
}
