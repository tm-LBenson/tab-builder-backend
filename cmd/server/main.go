package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tm-LBenson/tab-builder-backend/internal/db"
	"github.com/tm-LBenson/tab-builder-backend/internal/handlers"
	"github.com/tm-LBenson/tab-builder-backend/internal/middleware"

)

func main() {
	ctx := context.Background()
	store, err := db.NewStore(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	r := chi.NewRouter()
	r.Use(middleware.FirebaseAuth)

	song := handlers.SongHandler{Store: store}
	r.Route("/songs", song.Register)

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

