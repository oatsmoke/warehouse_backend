#!/bin/sh
set -e

echo "waiting for database..."
until pg_isready -h db -p 5432 -q -t 1; do
  sleep 1
done

echo "running Atlas migrations..."
atlas migrate apply --url "$POSTGRES_DSN" --dir "file://migrations"

echo "Add root user if not exists..."
psql "$POSTGRES_DSN" -f root_user.sql

echo "migrations applied. Starting application..."
exec "$@"