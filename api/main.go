package main

import (
	"log"
	"net/http"
	"os"

	"classdir/api/internal/auth"
	"classdir/api/internal/db"
	"classdir/api/internal/hub"
	"classdir/api/internal/presentation"
	"classdir/api/internal/shared/cfg"
)

func main() {
	pool := db.InitDB()

	port := os.Getenv(cfg.EnvPort)
	if port == "" {
		port = cfg.DefaultPort
	}

	mux := http.NewServeMux()

	auth.RegisterRoutes(mux)

	store := presentation.NewStore(pool)

	api := http.NewServeMux()
	presentation.RegisterRoutes(api, store)

	h := hub.NewHub(store)
	mux.Handle("GET /ws/v1", auth.AuthMiddleware(hub.WSHandler(h)))

	mux.Handle("/api/v1/", auth.AuthMiddleware(api))

	log.Printf("API server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
