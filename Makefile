ENV_FILE = ./docker/.env
-include $(ENV_FILE)

GOLANGCI_LINT_PATH = ./.golangci.yaml
COMPOSE_DEV=docker/docker-compose.dev.yml
COMPOSE_PROD=docker/docker-compose.prod.yml
DOCKER_COMPOSE ?= docker compose

MIGRATIONS_DIR = db/migrations
DB_URL = "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5434/$(POSTGRES_DB)?sslmode=disable"

WIRE_ROOT := internal/wire
WIRE_DIRS := $(patsubst %/,%,$(wildcard $(WIRE_ROOT)/*/))
WIRE_TARGETS := $(addsuffix .wire,$(WIRE_DIRS))

.PHONY: lint
lint: 
	@golangci-lint run --config $(GOLANGCI_LINT_PATH) --fix

# docker

.PHONY: docker-start-dev
docker-start-dev:
	@$(DOCKER_COMPOSE) -f $(COMPOSE_DEV) up -d

.PHONY: docker-build-dev
docker-build-dev:
	@$(DOCKER_COMPOSE) -f $(COMPOSE_DEV) build

.PHONY: docker-stop-dev
docker-stop-dev:
	@$(DOCKER_COMPOSE) -f $(COMPOSE_DEV) stop

.PHONY: docker-clean-dev
docker-clean-dev:
	@$(DOCKER_COMPOSE) -f $(COMPOSE_DEV) down

.PHONY: docker-start-prod
docker-start-prod:
	@$(DOCKER_COMPOSE) -f $(COMPOSE_PROD) up -d

.PHONY: docker-build-prod
docker-build-prod:
	@$(DOCKER_COMPOSE) -f $(COMPOSE_PROD) build

.PHONY: docker-stop-prod
docker-stop-prod:
	@$(DOCKER_COMPOSE) -f $(COMPOSE_PROD) stop

.PHONY: docker-clean-prod
docker-clean-prod:
	@$(DOCKER_COMPOSE) -f $(COMPOSE_PROD) down

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

EASYJSON_SOURCES := $(wildcard */*/dto/dto.go */*/*/dto/dto.go)

.PHONY: generate-easyjson
generate-easyjson:
	$(if $(EASYJSON_SOURCES),easyjson -all $(EASYJSON_SOURCES),@echo "No easyjson sources found")

# seed

.PHONY: seed
seed:
	docker exec -it proftwist-roadmap-service sh -c "./seed"

# wire

.PHONY: wire $(WIRE_TARGETS)
wire: $(WIRE_TARGETS)

$(WIRE_TARGETS):
	@echo "Running wire in $(@:.wire=)"
	@cd "$(@:.wire=)" && wire
