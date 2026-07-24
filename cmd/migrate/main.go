package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"clinic-vet-api/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	databaseURL, err := config.ResolveDatabaseURL()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Running migrations...")

	if err := waitForDatabase(databaseURL, 30*time.Second); err != nil {
		log.Fatalf("database not ready: %v", err)
	}

	if err := runMigrations(databaseURL); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	fmt.Println("Migrations completed successfully!")
}

func runMigrations(databaseURL string) error {
	m, err := migrate.New(
		"file://database/migrations",
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}

func waitForDatabase(databaseURL string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for database")
		default:
			pool, err := pgxpool.New(ctx, databaseURL)
			if err == nil {
				defer pool.Close()
				if err := pool.Ping(ctx); err == nil {
					return nil
				}
			}
			time.Sleep(2 * time.Second)
		}
	}
}
