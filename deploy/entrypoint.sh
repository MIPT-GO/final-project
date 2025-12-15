#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

if [ "${RUN_MIGRATIONS:-true}" = "true" ]; then
  readonly GOOSE_DRIVER=postgres
  readonly GOOSE_DBSTRING="${GOOSE_DRIVER}://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_DATABASE}?sslmode=disable"

  echo "Running database migrations..."
  if goose -dir /app/migrations/ "$GOOSE_DRIVER" "$GOOSE_DBSTRING" up; then
    echo "Migrations completed successfully!"
  else
    echo "Failed to run migrations!"
    exit 1
  fi
else
  echo "Migrations are disabled (RUN_MIGRATIONS=false), skipping..."
fi

echo "Starting application..."
exec "$@"
