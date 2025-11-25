ENV_FILE = ./docker/.env
include $(ENV_FILE)

GOLANGCI_LINT_PATH = ./.golangci.yaml
COMPOSE_DEV=docker/docker-compose.dev.yml
COMPOSE_PROD=docker/docker-compose.prod.yml

MIGRATIONS_DIR = db/migrations
DB_URL = "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable"

WIRE_ROOT := internal/wire
WIRE_DIRS := $(shell find $(WIRE_ROOT) -maxdepth 1 -mindepth 1 -type d)

.PHONY: lint
lint: 
	@golangci-lint run --config $(GOLANGCI_LINT_PATH) --fix

# docker

.PHONY: docker-start-dev
docker-start-dev:
	@docker compose -f $(COMPOSE_DEV) up -d

.PHONY: docker-build-dev
docker-build-dev:
	@docker compose -f $(COMPOSE_DEV) build

.PHONY: docker-stop-dev
docker-stop-dev:
	@docker compose -f $(COMPOSE_DEV) stop

.PHONY: docker-clean-dev
docker-clean-dev:
	@docker compose -f $(COMPOSE_DEV) down

.PHONY: docker-start-prod
docker-start-prod:
	@docker compose -f $(COMPOSE_PROD) up -d

.PHONY: docker-build-prod
docker-build-prod:
	@docker compose -f $(COMPOSE_PROD) build

.PHONY: docker-stop-prod
docker-stop-prod:
	@docker compose -f $(COMPOSE_PROD) stop

.PHONY: docker-clean-prod
docker-clean-prod:
	@docker compose -f $(COMPOSE_PROD) down

# migrations

.PHONY: migrate-create
migrate-create:
	@goose -dir $(MIGRATIONS_DIR) create $(name) sql

.PHONY: migrate-up
migrate-up:
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

.PHONY: migrate-down
migrate-down:
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

.PHONY: migrate-down-to-zero
migrate-down-to-zero:
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down-to 0

.PHONY: migrate-status
migrate-status:
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" status

.PHONY: migrate-version
migrate-version:
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" version

.PHONY: migrate-reset
migrate-reset:
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" reset

.PHONY: migrate-fix
migrate-fix:
	@goose -dir $(MIGRATIONS_DIR) fix

# easyjson

.PHONY: generate-easyjson
generate-easyjson:
	easyjson -all */*/dto/dto.go */*/*/dto/dto.go 

# seed

.PHONY: seed
seed:
	docker exec -it proftwist-roadmap-service sh -c "./seed";

# wire

.PHONY: wire
wire:
	@for dir in $(WIRE_DIRS); do \
		echo "Running wire in $$dir"; \
		( cd "$$dir" && wire ); \
	done
