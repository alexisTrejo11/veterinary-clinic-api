package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// ResolveDatabaseURL builds a PostgreSQL connection string for pgx / migrate.
//
// Preferred (.env):
//
//	DATABASE_URL=jdbc:postgresql://host:5432/dbname?sslmode=require
//	DATABASE_USER=postgres
//	DATABASE_PASSWORD=secret
//
// Also accepts a full URL with embedded credentials (legacy).
func ResolveDatabaseURL() (string, error) {
	raw := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if raw == "" {
		return "", fmt.Errorf("DATABASE_URL is required")
	}

	user := strings.TrimSpace(os.Getenv("DATABASE_USER"))
	if user != "" {
		password := os.Getenv("DATABASE_PASSWORD")
		return buildDatabaseURL(raw, user, password)
	}

	raw = normalizeDatabaseURL(raw)
	if !strings.Contains(raw, "@") {
		return "", fmt.Errorf("DATABASE_USER is required when DATABASE_URL has no credentials")
	}
	return raw, nil
}

func buildDatabaseURL(endpoint, user, password string) (string, error) {
	endpoint = normalizeDatabaseURL(endpoint)

	if !strings.Contains(endpoint, "://") {
		endpoint = "postgresql://" + endpoint
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("invalid DATABASE_URL: %w", err)
	}

	switch u.Scheme {
	case "postgres", "postgresql":
	default:
		return "", fmt.Errorf("unsupported DATABASE_URL scheme %q (use postgresql://)", u.Scheme)
	}

	if u.User != nil && u.User.Username() != "" {
		// Already a full URL (e.g. env left unchanged after entrypoint); use as-is.
		return endpoint, nil
	}

	u.User = url.UserPassword(user, password)

	if u.Query().Get("sslmode") == "" && strings.Contains(u.Host, "rds.amazonaws.com") {
		q := u.Query()
		q.Set("sslmode", "require")
		u.RawQuery = q.Encode()
	}

	return u.String(), nil
}

func normalizeDatabaseURL(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "jdbc:")
	return raw
}
