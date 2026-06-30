package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"classdir/api/internal/shared/cfg"
)

func InitDB() *pgxpool.Pool {
	dbURL := os.Getenv(cfg.EnvDatabaseURL)
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.DbTimeout)
	defer cancel()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to create pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	log.Println("Connected to database")
	return pool
}
