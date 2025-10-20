#!/bin/sh
set -e

echo "waiting for database..."
until nc -z db 5432; do
  sleep 1
done

echo "running Atlas migrations..."
atlas migrate apply --url "$POSTGRES_DSN" --dir "file://migrations"

echo "migrations applied. Starting application..."
exec "$@"