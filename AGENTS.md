# AGENTS.md — Proftwist Backend

## Overview

Go monorepo of microservices. Each service is a standalone binary in `cmd/{name}/main.go`.  
No ORM — raw SQL (PostgreSQL), MongoDB driver, Redis, MinIO/S3.  
Dependency injection via `google/wire`. Config via `spf13/viper` (YAML + env overrides).

---

## Architecture

### Service structure

```
services/{name}/
├── adapter.go          — interfaces for external adapters (Kafka, etc.)
├── adapter/            — adapter implementations
├── delivery.go          — HTTP, WS, gRPC interfaces
├── delivery/http/      — HTTP handlers + route registration
├── delivery/ws/        — WebSocket handlers
├── delivery/grpc/      — gRPC handlers
├── dto/                — request/response DTOs + mappers
├── repository.go       — repository interface
├── repository/         — repository implementation (Postgres, Mongo, Redis, etc.)
├── usecase.go          — usecase interface
└── usecase/            — usecase implementation (business logic)
```

### Internal packages

```
internal/
├── entities/           — domain models (plain structs, no methods)
│   ├── errs/           — sentinel errors: ErrNotFound, ErrForbidden, etc.
├── infrastructure/
│   ├── broker/         — Kafka producer/consumer
│   ├── client/         — gRPC clients for inter-service communication
│   └── db/             — DB connection setup (postgres, mongo, redis, aws)
├── server/
│   ├── grpc/           — shared gRPC server
│   ├── http/           — shared HTTP server (gorilla/mux)
│   ├── middleware/     — auth, cors, logging middleware
│   └── ws/             — WebSocket server
├── utils/              — helpers (JSON response, etc.)
├── wire/               — wire injection per service
└── worker/             — Kafka consumer workers
```

### Services (7)

| Service | DB | Ports |
|---------|----|-------|
| `cmd/auth/` | PostgreSQL, Redis, MinIO | auth, VK OAuth, file upload |
| `cmd/chat/` | PostgreSQL | group/direct chats, messages, threads |
| `cmd/category/` | PostgreSQL | roadmap categories |
| `cmd/friend/` | PostgreSQL | friend requests, friends |
| `cmd/roadmap/` | PostgreSQL + MongoDB | roadmap metadata (PG) + graph data (Mongo) |
| `cmd/roadmapinfo/` | PostgreSQL | roadmap info CRUD, search, subscriptions |
| `cmd/moderation/` | — | content moderation (AI) |

### Polyglot persistence

| Data | Store |
|------|-------|
| User accounts, auth | PostgreSQL |
| VK OAuth identities | PostgreSQL |
| Categories | PostgreSQL |
| Roadmap metadata (name, author, visibility) | PostgreSQL |
| Chat messages, group/direct chats | PostgreSQL |
| Friend requests, friendships | PostgreSQL |
| Roadmap graph (nodes, edges) | MongoDB |
| User progress on nodes | MongoDB |
| JWT blacklist, OAuth state | Redis |
| File uploads (avatars) | MinIO / S3 |

---

## Code Conventions

### Entities (`internal/entities/`)

Plain structs with `uuid.UUID` PKs, `time.Time` timestamps, no tags except `bson`/`json` where needed.

```go
type Message struct {
    ID           uuid.UUID
    ChatID       uuid.UUID
    UserID       uuid.UUID
    Content      string
    ThreadRootID *uuid.UUID    // nullable pointer for optional FKs
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

- Use `*uuid.UUID` for nullable foreign keys
- Use `*string` for optional text fields
- Always include `CreatedAt` / `UpdatedAt`

### DTOs (`services/{name}/dto/`)

JSON tags with `omitempty` for optional fields, `validate:"required"` for validation.

```go
type ChatMessageResponseDTO struct {
    ID           uuid.UUID         `json:"id"`
    ThreadRootID *uuid.UUID        `json:"thread_root_id,omitempty"`
    ReplyCount   int               `json:"reply_count"`
}
```

> **IMPORTANT:** After adding fields to any DTO struct that has a `*_easyjson.go` file, run `make generate-easyjson` to regenerate the custom marshalers. Otherwise the new fields are silently dropped during JSON serialization.

### Repository pattern

Interface in `services/{name}/repository.go`, impl in `services/{name}/repository/`.

- Use `database/sql` with `github.com/lib/pq` for PostgreSQL
- Use `go.mongodb.org/mongo-driver` for MongoDB
- Use `github.com/redis/go-redis/v9` for Redis
- Use `github.com/minio/minio-go/v7` for S3

SQL queries are defined as package-level constants in `queries.go`.

### Usecase pattern

Interface in `services/{name}/usecase.go`, impl in `services/{name}/usecase/`.

- Receives DTOs, returns DTOs
- Orchestrates repo calls, external clients, and adapters
- Enriches responses with user data via gRPC auth client
- All methods accept `context.Context` as first param

### Error handling

Sentinel errors in `internal/entities/errs/errors.go`:

```go
var (
    ErrNotFound           = errors.New("not found")
    ErrForbidden          = errors.New("forbidden")
    ErrAlreadyExists      = errors.New("already exists")
    ErrUnauthorized       = errors.New("unauthorized")
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrBusinessLogic      = errors.New("business logic error")
)

// Check helpers: IsNotFoundError, IsForbiddenError, IsAlreadyExistsError, IsBusinessLogicError
```

Handlers check errors and map to HTTP status codes:
```go
if errs.IsNotFoundError(err) { statusCode = http.StatusNotFound }
if errs.IsForbiddenError(err) { statusCode = http.StatusForbidden }
```

### HTTP handlers

- `gorilla/mux` router
- Routes registered via `HttpRegistrar` interface
- Auth via `AuthMiddleware` (reads JWT from `jwt-token` cookie or `Authorization` header)
- Request context carries user ID via `utils.UserIDKey{}`
- Responses via `utils.JSONResponse(ctx, w, statusCode, body)` and `utils.JSONError(ctx, w, statusCode, msg)`
- Routes use `/api/v1/{resource}/...` prefix

### WebSocket

- Message envelope: `{"type": "string", "data": {...}, "timestamp": "RFC3339"}`
- Handler registration via `WsRegistrar`
- Outgoing types: `message_sent`, `typing_notification`, `user_joined`, `user_left`
- Incoming types: `send_message`, `typing`
- Events routed through Kafka (notification service → WS broadcast)

### gRPC

- Inter-service communication via gRPC clients in `internal/infrastructure/client/`
- Proto definitions in each client's `proto/` subdirectory
- Regenerate: `protoc --go_out=... --go-grpc_out=... proto/*.proto`

### Configuration

- `spf13/viper` reads `config.yml` + env overrides
- Config struct in `config/config.go`
- `.env` file in `docker/.env`
- Env vars map to config keys via `bindEnv()`

### Dependency injection (Wire)

- Wire sets in `internal/wire/{service}/`
- Run `make wire` after changing dependencies
- Each service's `main.go` calls `wire.Build(...)` from its wire package

---

## Common Commands

```bash
make docker-start-dev         # Start all infrastructure + services
make docker-build-dev         # Rebuild all service images
make docker-stop-dev          # Stop all containers
make docker-clean-dev         # Down all containers

# Single service (replace {name})
docker compose -f docker/docker-compose.dev.yml build --no-cache {name}-service
docker compose -f docker/docker-compose.dev.yml up -d {name}-service

make migrate-up               # Run pending DB migrations
make migrate-down             # Rollback last migration
make migrate-create name=xxx  # Create new migration file
make migrate-status           # Check migration status

make generate-easyjson        # Regenerate all easyjson marshalers
make wire                     # Regenerate wire injection code
make lint                     # Run golangci-lint
make seed                     # Seed database (inside roadmap container)
```

---

## Database Migrations

- Tool: `goose`
- Location: `db/migrations/`
- File naming: `YYYYMMDDHHMMSS_description.sql`
- Use `-- +goose Up` / `-- +goose Down` statement separators

```sql
-- +goose Up
ALTER TABLE group_chat_messages ADD COLUMN thread_root_id UUID REFERENCES group_chat_messages(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE group_chat_messages DROP COLUMN thread_root_id;
```

---

## CRUD Pattern Example

```go
// services/{name}/usecase/usecase.go
func (uc *ChatUsecase) GetGroupChatMessages(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, limit, offset int) (*dto.GetChatMessagesResponseDTO, error) {
    // 1. Authorize
    isMember, err := uc.repo.IsGroupChatMember(ctx, chatID, userID)
    if err != nil { return nil, fmt.Errorf("...: %w", err) }
    if !isMember { return nil, errs.ErrForbidden }

    // 2. Fetch from repo
    messages, err := uc.repo.GetGroupChatMessages(ctx, chatID, limit, offset)

    // 3. Enrich (user data via gRPC, reply counts, etc.)
    userData := uc.fetchUserData(ctx, userID, userIDs)
    replyCounts, _ := uc.repo.GetThreadReplyCounts(ctx, rootIDs)
    response := dto.GetChatMessagesResponseToDTO(messages, userData)
    // merge replyCounts into response...

    // 4. Return
    return &response, nil
}
```

---

## Key Libraries

| Library | Usage |
|---------|-------|
| `gorilla/mux` | HTTP routing |
| `gorilla/websocket` | WebSocket server |
| `lib/pq` | PostgreSQL driver |
| `go.mongodb.org/mongo-driver` | MongoDB driver |
| `go-redis/v9` | Redis client |
| `minio-go/v7` | MinIO/S3 client |
| `google/wire` | Dependency injection |
| `spf13/viper` | Configuration |
| `golang-jwt/jwt/v4` | JWT auth |
| `mailru/easyjson` | Fast JSON marshal/unmarshal (generated) |
| `golang/protobuf` | gRPC protobuf |
| `google.golang.org/grpc` | gRPC framework |
| `segmentio/kafka-go` | Kafka producer/consumer |
| `google/uuid` | UUID generation |
| `sirupsen/logrus` | Structured logging |
| `prometheus/client_golang` | Metrics |
