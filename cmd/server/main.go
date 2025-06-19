package main

import (
	"context"
	"log"
	"net/http"
	"os"

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	r := chi.NewRouter()
	r.Use(middleware.FirebaseAuth)

	song := handlers.SongHandler{Store: store}
	r.Route("/songs", song.Register)

	log.Printf("listening on :%s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
