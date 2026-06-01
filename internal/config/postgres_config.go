package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseConfig struct {
	URL             string        `json:"url"`
	MaxOpenConns    int           `json:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time"`
}

var pgPool *pgxpool.Pool

func PostgresConn(dbURL string) *pgx.Conn {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn
}

func CreatePgxPool(dbURL string) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse database config: %v\n", err)
		os.Exit(1)
	}

	// Configurar el pool de conexiones
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	// Verificar la conexión
	if err := pool.Ping(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	// Store pool globally for later cleanup
	pgPool = pool

	fmt.Println("Successfully connected to PostgreSQL database with pgx pool!")
	return pool
}

// ClosePgxPool closes the PostgreSQL connection pool gracefully
func ClosePgxPool() {
	if pgPool != nil {
		pgPool.Close()
		pgPool = nil
		fmt.Println("PostgreSQL connection pool closed successfully")
	}
}

func loadDatabaseConfig(config *DatabaseConfig) error {
	dbURL, err := ResolveDatabaseURL()
	if err != nil {
		return err
	}
	config.URL = dbURL

	config.MaxOpenConns, err = parseIntWithDefault("DATABASE_MAX_OPEN_CONNS", 25)
	if err != nil {
		return fmt.Errorf("invalid DATABASE_MAX_OPEN_CONNS: %w", err)
	}

	config.MaxIdleConns, err = parseIntWithDefault("DATABASE_MAX_IDLE_CONNS", 5)
	if err != nil {
		return fmt.Errorf("invalid DATABASE_MAX_IDLE_CONNS: %w", err)
	}

	config.ConnMaxLifetime, err = parseDuration("DATABASE_CONN_MAX_LIFETIME", "1h")
	if err != nil {
		return fmt.Errorf("invalid DATABASE_CONN_MAX_LIFETIME: %w", err)
	}

	config.ConnMaxIdleTime, err = parseDuration("DATABASE_CONN_MAX_IDLE_TIME", "30m")
	if err != nil {
		return fmt.Errorf("invalid DATABASE_CONN_MAX_IDLE_TIME: %w", err)
	}

	return nil
}
