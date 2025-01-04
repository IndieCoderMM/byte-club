.PHONY: all
all:
	@echo "Specify a command:\n migration, migrate, migrate-down, migrate-status, migrate-validate"

.PHONY: migration
migration:
	@if [ "$(filter-out $@,$(MAKECMDGOALS))" = "" ]; then \
		echo "Error: Require migration name, e.g., 'make migration new_migration'"; \
		exit 1; \
	fi
	@goose create $(filter-out $@,$(MAKECMDGOALS)) sql

%:
	@:

.PHONY: migrate
migrate:
	@goose up

.PHONY: migrate-down
migrate-down:
	@goose down

.PHONY: migrate-status
migrate-status:
	@goose status

.PHONY: migrate-validate
migrate-validate:
	@goose validate
