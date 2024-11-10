# Include environment variables from .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

# Construct connection URLs
POSTGRES_URL=postgres://$(PG_USERNAME):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DATABASE)?sslmode=disable&search_path=$(PG_SCHEMA)
CASSANDRA_URL=cassandra://$(CASSANDRA_HOST):$(CASSANDRA_PORT)/$(CASSANDRA_KEYSPACE)?username=$(CASSANDRA_USERNAME)&password=$(CASSANDRA_PASSWORD)

# Migration paths
BASE_MIGRATION_PATH=internals/database/migrations
POSTGRES_MIGRATION_PATH=$(BASE_MIGRATION_PATH)/postgres
CASSANDRA_MIGRATION_PATH=$(BASE_MIGRATION_PATH)/cassandra
MIGRATE_CMD=migrate

# Build the application
all: build test

build:
	@echo "Building..."
	@go build -o temp/main cmd/server/main.go

# Run the application
run:
	@go run cmd/main.go

# Docker commands for database containers
docker-run:
	@if docker compose up -d 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up -d; \
	fi

docker-pause:
	@if docker compose pause 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose pause; \
	fi

docker-unpause:
	@if docker compose unpause 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose unpause; \
	fi

docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test commands
test:
	@echo "Testing..."
	@go test ./... -v

itest:
	@echo "Running integration tests..."
	@go test ./internals/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f tmp/main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
		air; \
		echo "Watching...";\
	else \
		read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/air-verse/air@latest; \
			air; \
			echo "Watching...";\
		else \
			echo "You chose not to install air. Exiting..."; \
			exit 1; \
		fi; \
	fi

# Database migration tools installation
.PHONY: install-tools
install-tools:
	@echo "Installing migration tools..."
	@go install -tags 'cassandra postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Ensure migration directories exist
.PHONY: ensure-dirs
ensure-dirs:
	@mkdir -p $(POSTGRES_MIGRATION_PATH)
	@mkdir -p $(CASSANDRA_MIGRATION_PATH)

# Print current database configuration
.PHONY: print-db-config
print-db-config:
	@echo "PostgreSQL Configuration:"
	@echo "  Host: $(PG_HOST)"
	@echo "  Port: $(PG_PORT)"
	@echo "  Database: $(PG_DATABASE)"
	@echo "  User: $(PG_USER)"
	@echo "  SSL Mode: $(PG_SSL_MODE)"
	@echo "\nCassandra Configuration:"
	@echo "  Host: $(CASSANDRA_HOST)"
	@echo "  Port: $(CASSANDRA_PORT)"
	@echo "  Keyspace: $(CASSANDRA_KEYSPACE)"
	@echo "  User: $(CASSANDRA_USER)"

# PostgreSQL migration commands
create-postgres-migration: ensure-dirs
	@read -p "Enter migration name: " name; \
	$(MIGRATE_CMD) create -ext sql -dir $(POSTGRES_MIGRATION_PATH) -seq $$name

migrate-postgres-up: ensure-dirs
	$(MIGRATE_CMD) -path $(POSTGRES_MIGRATION_PATH) -database "$(POSTGRES_URL)" up

migrate-postgres-down:
	$(MIGRATE_CMD) -path $(POSTGRES_MIGRATION_PATH) -database "$(POSTGRES_URL)" down 1

migrate-postgres-reset:
	$(MIGRATE_CMD) -path $(POSTGRES_MIGRATION_PATH) -database "$(POSTGRES_URL)" drop -f
	$(MIGRATE_CMD) -path $(POSTGRES_MIGRATION_PATH) -database "$(POSTGRES_URL)" up

# Cassandra migration commands
create-cassandra-migration: ensure-dirs
	@read -p "Enter migration name: " name; \
	$(MIGRATE_CMD) create -ext cql -dir $(CASSANDRA_MIGRATION_PATH) -seq $$name

migrate-cassandra-up: ensure-dirs
	$(MIGRATE_CMD) -path $(CASSANDRA_MIGRATION_PATH) -database "$(CASSANDRA_URL)" up

migrate-cassandra-down:
	$(MIGRATE_CMD) -path $(CASSANDRA_MIGRATION_PATH) -database "$(CASSANDRA_URL)" down 1

migrate-cassandra-reset:
	$(MIGRATE_CMD) -path $(CASSANDRA_MIGRATION_PATH) -database "$(CASSANDRA_URL)" drop -f
	$(MIGRATE_CMD) -path $(CASSANDRA_MIGRATION_PATH) -database "$(CASSANDRA_URL)" up

# Migration status
migrate-status:
	@echo "PostgreSQL migration status:"
	$(MIGRATE_CMD) -path $(POSTGRES_MIGRATION_PATH) -database "$(POSTGRES_URL)" version
	@echo "\nCassandra migration status:"
	$(MIGRATE_CMD) -path $(CASSANDRA_MIGRATION_PATH) -database "$(CASSANDRA_URL)" version

.PHONY: all build run test clean watch docker-run docker-down docker-pause docker-unpause \
	itest create-postgres-migration migrate-postgres-up migrate-postgres-down \
	migrate-postgres-reset create-cassandra-migration migrate-cassandra-up \
	migrate-cassandra-down migrate-cassandra-reset migrate-status print-db-config