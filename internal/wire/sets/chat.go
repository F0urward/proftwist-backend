package sets

import (
	"github.com/google/wire"

	chatAdapter "github.com/F0urward/proftwist-backend/services/chat/adapter"
	chatGRPCHandlers "github.com/F0urward/proftwist-backend/services/chat/delivery/grpc"
	chatHTTPHandlers "github.com/F0urward/proftwist-backend/services/chat/delivery/http"
	chatWSHandlers "github.com/F0urward/proftwist-backend/services/chat/delivery/ws"
	chatRepository "github.com/F0urward/proftwist-backend/services/chat/repository"
	chatUsecase "github.com/F0urward/proftwist-backend/services/chat/usecase"
)

var ChatSet = wire.NewSet(
	chatRepository.NewChatPostgresRepository,
	chatUsecase.NewChatUsecase,
	chatHTTPHandlers.NewChatHandler,
	chatWSHandlers.NewChatWSHandlers,
	chatGRPCHandlers.NewChatServer,
	chatAdapter.NewWSNotifier,
)
