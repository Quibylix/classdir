package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	initDB()

	port := os.Getenv(envPort)
	if port == "" {
		port = defaultPort
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/auth/login", loginHandler)
	mux.HandleFunc("POST /api/v1/auth/logout", logoutHandler)

	presStore := newPresentationStore(db)

	api := http.NewServeMux()
	api.HandleFunc("POST /api/v1/presentation", createPresentationHandler(presStore))

	mux.Handle("/api/v1/", authMiddleware(api))

	log.Printf("API server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
