NAME ?= migration
ENV ?= local

migrate_diff:
	atlas migrate diff $(NAME) --env local

migrate_apply:
	atlas migrate apply --env $(ENV)

init_env:
	$env:POSTGRES_DSN = "postgres://root:password@localhost:5432/wh?sslmode=disable"
	$env:TEST_POSTGRES_DSN = "postgres://test:password@localhost:55432/test?sslmode=disable"