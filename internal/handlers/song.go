package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tm-LBenson/tab-builder-backend/internal/db"
	"github.com/tm-LBenson/tab-builder-backend/internal/middleware"
)

type SongHandler struct{ Store *db.Store }


func (h SongHandler) listMine(w http.ResponseWriter, r *http.Request) {
	uid := middleware.UID(r)
	rows, err := h.Store.ListUserSongs(r.Context(), uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(rows)
}

func (h SongHandler) listPublic(w http.ResponseWriter, r *http.Request) {
	lim := 10
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 50 {
			lim = n
		}
	}
	rows, err := h.Store.ListPublicSongs(r.Context(), lim)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(rows)
}


func (h SongHandler) create(w http.ResponseWriter, r *http.Request) {
	uid := middleware.UID(r)
	var in db.SongIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	out, err := h.Store.CreateSong(r.Context(), uid, in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(out)
}

func (h SongHandler) get(w http.ResponseWriter, r *http.Request) {
	uid := middleware.UID(r)
	id  := chi.URLParam(r, "id")

	song, err := h.Store.GetSong(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if song.OwnerUID != uid && !song.IsPublic {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	json.NewEncoder(w).Encode(song)
}

func (h SongHandler) update(w http.ResponseWriter, r *http.Request) {
	uid := middleware.UID(r)
	id  := chi.URLParam(r, "id")

	orig, err := h.Store.GetSong(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if orig.OwnerUID != uid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var in db.SongIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	out, err := h.Store.UpdateSong(r.Context(), id, in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(out)
}

func (h SongHandler) delete(w http.ResponseWriter, r *http.Request) {
	uid := middleware.UID(r)
	id  := chi.URLParam(r, "id")

	song, err := h.Store.GetSong(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if song.OwnerUID != uid {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if err := h.Store.DeleteSong(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h SongHandler) clone(w http.ResponseWriter, r *http.Request) {
	uid := middleware.UID(r)
	id  := chi.URLParam(r, "id")

	out, err := h.Store.CloneSong(r.Context(), id, uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(out)
}


func (h SongHandler) Register(r chi.Router) {
	// Public feed (no auth)
	r.Get("/public", h.listPublic)

	// Authenticated endpoints
	r.Get("/", h.listMine)
	r.Post("/", h.create)

	r.Route("/{id}", func(sr chi.Router) {
		sr.Get("/", h.get)
		sr.Put("/", h.update)
		sr.Delete("/", h.delete)
		sr.Post("/clone", h.clone)
	})
}
