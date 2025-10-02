ENV_FILE = ./docker/.env
include $(ENV_FILE)

GOLANGCI_LINT_PATH = ./.golangci.yaml
DOCKER_COMPOSE_PATH=docker/docker-compose.yml

MIGRATIONS_DIR = db/migrations
DB_URL = "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable"

.PHONY: lint

lint: 
	@golangci-lint run --config $(GOLANGCI_LINT_PATH) --fix

# Docker

.PHONY: docker-build

docker-build:
	@docker compose -f $(DOCKER_COMPOSE_PATH) build

.PHONY: docker-start

docker-start:
	@docker compose -f $(DOCKER_COMPOSE_PATH) up -d

.PHONY: docker-stop

docker-stop:
	@docker compose -f $(DOCKER_COMPOSE_PATH) stop

.PHONY: docker-clean

docker-clean:
	@docker compose -f $(DOCKER_COMPOSE_PATH) down

# Migrations

.PHONY: migrate-create
migrate-create:
	@migrate create -ext sql -dir $(MIGRATIONS_DIR) $(name)

.PHONY: migrate-up
migrate-up:
	migrate -database $(DB_URL) -path $(MIGRATIONS_DIR) up

.PHONY: migrate-down
migrate-down:
	@migrate -database $(DB_URL) -path $(MIGRATIONS_DIR) down 1

.PHONY: migrate-force
migrate-force:
	@migrate -database $(DB_URL) -path $(MIGRATIONS_DIR) force $(version)

