# --- CONFIGURATION VARIABLES ---
# DB & Tools
SQLC_VERSION := v1.30.0
GOOSE_VERSION := v3.21.1
MIGRATIONS_DIR := migrations

APP_DIR := app
APP_OUTPUT_DIR := app/bin
WEB_DIR := web
WEB_OUTPUT_DIR := web/dist

# Include .env file if it exists
-include .env.defaults
-include .env

# Set default DB_URL if not defined in .env
DB_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)

.PHONY: all help clean \
		install-tools sqlc migrate-up migrate-down migrate-down-all migrate-create migrate-status \
		run-migration test-migration \
		env-setup build-app test-app build-web test-web

# Default target
help:
	@echo "Available targets:"
	@echo "  --- BUILD & DEPLOY ---"
	@echo "  build-app    	  - Build all Go services from cmd/ (OS: $(GOOS), Arch: $(GOARCH))"
	@echo "  build-web        - Build React frontend (installs deps + build)"
	@echo "  clean            - Remove build artifacts (bin/ and dist/)"
	@echo "  all              - Clean and build everything"
	@echo ""
	@echo "  --- DATABASE & TOOLS ---"
	@echo "  install-tools    - Install SQLC and goose tools"
	@echo "  sqlc             - Generate Go code from SQL queries using SQLC"
	@echo "  migrate-up       - Run all pending migrations"
	@echo "  migrate-down     - Rollback the last migration"
	@echo "  migrate-down-all - Rollback all applied migrations (down to version 0)"
	@echo "  migrate-create   - Create a new migration file (usage: make migrate-create name=my_migration)"
	@echo "  migrate-status   - Show migration status"
	@echo "  run-migration    - Run migration sequence"
	@echo "  test-migration	  - Run CI test migration sequence"
	@echo "  env-setup        - Create a .env file from template"

# --- META TARGETS ---
all: clean build-app build-web

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(APP_OUTPUT_DIR)
	@rm -rf $(WEB_OUTPUT_DIR)
	@echo "Done."

# --- BUILD TARGETS ---

# Build Generic Go Binaries (scans cmd/ directory)
build-app:
	@echo "=== Building Backend Services (OS: $(GOOS), Arch: $(GOARCH)) ==="
	# 1. Clean and create base output directory
	@rm -rf $(APP_OUTPUT_DIR)
	@mkdir -p $(APP_OUTPUT_DIR)

	# Build binaries
	@for app_dir in ${APP_DIR}/cmd/*; do \
		if [ -d "$$app_dir" ] && [ -f "$$app_dir/main.go" ]; then \
			app_name=$$(basename "$$app_dir"); \
			echo " > Building $$app_name..."; \
			go build -C "$(APP_DIR)" -v -o "bin/$$app_name" "./cmd/$$app_name/main.go" || exit 1; \
		fi \
	done
	@echo "=== Backend Build Complete (Check $(APP_OUTPUT_DIR)/) ==="

test-app:
	@echo "=== Running Go tests in $(APP_DIR)/ ==="
	@go test -C "$(APP_DIR)" -coverprofile=coverage.out -covermode=atomic $$(go list -C "$(APP_DIR)" ./... | grep -v '/gen$$') 2>&1 || true
	@echo "=== Go tests finished ==="
	@echo "=== Generating coverage report ==="
	@if [ -f $(APP_DIR)/coverage.out ]; then \
		go tool cover -html=$(APP_DIR)/coverage.out -o $(APP_DIR)/coverage.html 2>/dev/null && \
		echo "Coverage report saved to $(APP_DIR)/coverage.html" || \
		echo "Coverage data saved to $(APP_DIR)/coverage.out"; \
	else \
		echo "No coverage data generated"; \
	fi

build-web:
	@echo "=== Building Web Frontend ==="
	@if [ ! -d "$(WEB_DIR)" ]; then echo "Error: $(WEB_DIR) directory not found"; exit 1; fi
	@echo " > Installing dependencies..."
	@cd $(WEB_DIR) && npm ci
	@echo " > Running build..."
	@cd $(WEB_DIR) && npm run build
	@if [ ! -d "$(WEB_OUTPUT_DIR)" ]; then echo "Error: dist folder missing after build"; exit 1; fi
	@echo "=== Web Build Complete ==="

test-web:
	@echo "=== Running Web Frontend Tests (with fresh deps) ==="
	@if [ ! -d "$(WEB_DIR)" ]; then echo "Error: $(WEB_DIR) directory not found"; exit 1; fi
	@cd $(WEB_DIR) && npm ci
	@cd $(WEB_DIR) && npm run test:coverage
	@echo "=== Web Tests Finished ==="


# --- DATABASE & TOOLS TARGETS ---

install-tools:
	@echo "Installing tools..."
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION)
	go install github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)
	@echo "Tools installed successfully!"

sqlc:
	@echo "Generating Go code from SQL queries..."
	sqlc generate -f app/sqlc.yml
	@echo "Code generation completed!"

migrate-up:
	@echo "Running migrations (goose up)..."
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up
	@echo "Migrations completed!"

migrate-down:
	@echo "Rolling back the last migration (goose down 1)..."
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down 1
	@echo "Rollback completed!"

migrate-down-all:
	@echo "Rolling back all migrations (goose down-to 0)..."
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down-to 0
	@echo "All migrations rolled back!"

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: Migration name not provided. Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	@echo "Creating new migration (sequential): $(name)"
	@goose -dir $(MIGRATIONS_DIR) -s create $(name) sql
	@echo "Migration file created!"

migrate-status:
	@echo "Migration status:"
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" status

run-migration:
	@echo "=== Running migration sequence ==="
	$(MAKE) migrate-status
	$(MAKE) migrate-up || exit 1
	$(MAKE) migrate-status
	@echo "=== Migration sequence completed ==="

test-migration:
	@echo "=== Running CI test migration sequence ==="
	$(MAKE) migrate-status
	$(MAKE) migrate-up || exit 1
	$(MAKE) migrate-down-all || exit 1
	$(MAKE) migrate-up || exit 1
	$(MAKE) migrate-down || exit 1
	$(MAKE) migrate-up || exit 1
	$(MAKE) migrate-status
	@echo "=== Migration test sequence completed successfully ==="

env-setup:
	@if [ -f .env ]; then \
		echo ".env file already exists. Remove it first if you want to create a new one."; \
	else \
		cp .env.defaults .env; \
		echo ".env file created from template. Please update it with your credentials."; \
	fi