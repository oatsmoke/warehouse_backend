NAME ?= migration
ENV ?= local

migrate_diff:
	atlas migrate diff $(NAME) --env $(ENV)

migrate_apply:
	atlas migrate apply --env $(ENV)