NAME ?= migration
SCHEMA := file://schema/init_up.sql
DEV_URL := docker://postgres/17/dev?search_path=public
MIGRATIONS_DIR := file://migrations

migrate:
	atlas migrate diff "$(NAME)" --to "$(SCHEMA)" --dev-url "$(DEV_URL)" --dir "$(MIGRATIONS_DIR)"
