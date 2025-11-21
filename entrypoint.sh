#!/bin/sh
set -e

echo "waiting for database..."
until pg_isready -h wh -p 5432 -q -t 1; do
  sleep 1
done

echo "running atlas migrations..."
atlas migrate apply --env local

echo "add root user if not exists..."
psql "$POSTGRES_DSN" -f root_user.sql

echo "migrations applied. Starting application..."
exec "$@"