package chat

import (
	grpcServer "github.com/F0urward/proftwist-backend/internal/server/grpc"
	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
	wsServerHTTPHandlers "github.com/F0urward/proftwist-backend/internal/server/ws/http"
	chat "github.com/F0urward/proftwist-backend/services/chat"
	chatHTTPHandlers "github.com/F0urward/proftwist-backend/services/chat/delivery/http"
)

func AllHttpRegistrars(
	chatHandlers chat.Handlers,
	wsHandlers *wsServerHTTPHandlers.WebSocketHandler,
) []httpServer.HttpRegistrar {
	chatRegistrar := chatHTTPHandlers.NewChatHttpRegistrar(chatHandlers)
	wsRegistrar := wsServerHTTPHandlers.NewWebSocketHttpRegistrar(wsHandlers)

	return []httpServer.HttpRegistrar{
		chatRegistrar,
		wsRegistrar,
	}
}

func AllGrpcRegistrars(
	chatGrpcRegistrar grpcServer.GrpcRegistrar,
) []grpcServer.GrpcRegistrar {
	return []grpcServer.GrpcRegistrar{
		chatGrpcRegistrar,
	}
}
