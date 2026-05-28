#!/bin/sh
set -e

echo "Running database migrations..."

export PGHOST="${DATABASE_HOST:-localhost}"
export PGPORT="${DATABASE_PORT:-5432}"
export PGDATABASE="${DATABASE_NAME:-vivery}"
export PGUSER="${DATABASE_USER:-vivery}"
export PGPASSWORD="${DATABASE_PASSWORD:-vivery}"
export PGSSLMODE="${DATABASE_SSLMODE:-disable}"

/app/tern migrate --migrations /app/migrations

echo "Starting vivery..."
exec /app/vivery
