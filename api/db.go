package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func initDB() {
	dbURL := os.Getenv(envDatabaseURL)
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	var err error
	dbTimeoutCtx, cancelDbTimeout := context.WithTimeout(context.Background(), dbTimeout)
	defer cancelDbTimeout()
	db, err = pgxpool.New(dbTimeoutCtx, dbURL)
	if err != nil {
		log.Fatalf("failed to create pool: %v", err)
	}

	if err := db.Ping(dbTimeoutCtx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	log.Println("Connected to database")
}
