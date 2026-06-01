#!/bin/sh
set -e

echo "Starting application entrypoint..."

if [ -f .env ]; then
	echo "Loading environment variables from .env file"
	set -a
	# shellcheck disable=SC1091
	. ./.env
	set +a
fi

if [ -z "${DATABASE_URL:-}" ]; then
	echo "DATABASE_URL is required" >&2
	exit 1
fi

# shellcheck disable=SC1091
. /scripts/parse-database-url.sh

echo "Database target: ${DB_USER}@${DB_HOST}:${DB_PORT}/${DB_NAME}"

if [ "${SKIP_MIGRATIONS}" = "true" ]; then
	echo "Skipping migrations (SKIP_MIGRATIONS=true)"
elif command -v migrate >/dev/null 2>&1; then
	echo "Running database migrations..."

	if [ "${SKIP_DB_WAIT}" != "true" ]; then
		echo "Waiting for PostgreSQL to be ready..."
		until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER"; do
			echo "Waiting for database connection..."
			sleep 2
		done
		echo "PostgreSQL is ready!"
	else
		echo "Skipping database readiness wait (SKIP_DB_WAIT=true)"
	fi

	echo "Executing migrations..."
	if migrate -path ./db/migrations -database "$RESOLVED_DATABASE_URL" -verbose up; then
		echo "Migrations completed successfully!"
	else
		echo "Migrations failed; continuing application startup..."
	fi
else
	echo "migrate command not found, skipping migrations"
fi

echo "Starting application..."
exec "$@"
