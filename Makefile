NAME ?= migration
SCHEMA := file://schema/schema.sql
DEV_URL := docker://postgres/17/dev
MIGRATIONS_DIR := file://migrations

migrate_diff:
	atlas migrate diff "$(NAME)" --to "$(SCHEMA)" --dev-url "$(DEV_URL)" --dir "$(MIGRATIONS_DIR)"

test_db_up:
	docker compose up -d