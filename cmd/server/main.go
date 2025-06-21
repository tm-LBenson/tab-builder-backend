package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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

	frontend := os.Getenv("FRONTEND_ORIGIN")
	if frontend == "" {
		frontend = "http://localhost:5173"
	}
	corsMw := cors.New(cors.Options{
		AllowedOrigins:   []string{frontend},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r := chi.NewRouter()
	r.Use(corsMw.Handler)
	r.Use(middleware.FirebaseAuth(store))

	song := handlers.SongHandler{Store: store}
	r.Route("/songs", song.Register)

	log.Printf("listening on :%s, allowing CORS from %s\n", port, frontend)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
