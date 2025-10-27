ENV_FILE = ./docker/.env
include $(ENV_FILE)

GOLANGCI_LINT_PATH = ./.golangci.yaml
DOCKER_COMPOSE_PATH=docker/docker-compose.yml

MIGRATIONS_DIR = db/migrations
DB_URL = "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5434/$(POSTGRES_DB)?sslmode=disable"

.PHONY: lint

lint: 
	@golangci-lint run --config $(GOLANGCI_LINT_PATH) --fix

# docker

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

.PHONY: migrate-fix
migrate-fix:
	@goose -dir $(MIGRATIONS_DIR) fix

.PHONY: generate-easyjson
generate-easyjson:
	easyjson -all services/*/dto/dto.go

.PHONY: seed
seed:
	docker exec -it proftwist sh -c "./seed"; \
