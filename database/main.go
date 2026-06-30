package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

func main() {
	ctx := context.Background()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	if err := runMigrations(ctx, conn); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Migrations applied successfully")
}

func runMigrations(ctx context.Context, conn *pgx.Conn) error {
	_, err := conn.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY);`)
	if err != nil {
		return err
	}

	files, err := migrationFS.ReadDir("migrations")
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		var count int
		err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = $1", file.Name()).Scan(&count)
		if err != nil {
			return err
		}

		if count == 0 {
			fmt.Printf("Aplicando migración: %s\n", file.Name())

			content, err := migrationFS.ReadFile("migrations/" + file.Name())
			if err != nil {
				return err
			}

			_, err = conn.Exec(ctx, string(content))
			if err != nil {
				return fmt.Errorf("error en %s: %w", file.Name(), err)
			}

			_, err = conn.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", file.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}
