package chat

import (
	"github.com/google/wire"

	chatAdapter "github.com/F0urward/proftwist-backend/services/chat/adapter"
	chatGRPCHandlers "github.com/F0urward/proftwist-backend/services/chat/delivery/grpc"
	chatHTTPHandlers "github.com/F0urward/proftwist-backend/services/chat/delivery/http"
	chatRepository "github.com/F0urward/proftwist-backend/services/chat/repository"
	chatUsecase "github.com/F0urward/proftwist-backend/services/chat/usecase"

	authClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/authclient"
	db "github.com/F0urward/proftwist-backend/internal/infrastructure/db/postgres"

	wsServerHTTPHandlers "github.com/F0urward/proftwist-backend/internal/server/ws/http"
)

var ChatSet = wire.NewSet(
	chatRepository.NewChatPostgresRepository,
	chatUsecase.NewChatUsecase,
	chatHTTPHandlers.NewChatHandler,
	chatGRPCHandlers.NewChatServer,
	chatGRPCHandlers.NewChatGrpcRegistrar,
	wsServerHTTPHandlers.NewWebSocketHandler,
	chatAdapter.NewWSNotifier,
)

var ClientsSet = wire.NewSet(
	db.NewDatabase,
	authClient.NewAuthClient,
)
