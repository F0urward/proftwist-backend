//go:build wireinject
// +build wireinject

package chat

import (
	"github.com/google/wire"

	"github.com/F0urward/proftwist-backend/config"
	grpcServer "github.com/F0urward/proftwist-backend/internal/server/grpc"
	httpServer "github.com/F0urward/proftwist-backend/internal/server/http"
	authmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/auth"
	corsmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/cors"
	wsServer "github.com/F0urward/proftwist-backend/internal/server/ws"
	chatWSHandlers "github.com/F0urward/proftwist-backend/services/chat/delivery/ws"
)

func InitializeChatWsServer(cfg *config.Config) *wsServer.WsServer {
	wire.Build(
		wsServer.New,
	)
	return &wsServer.WsServer{}
}

func IntitializeChatWsRegistrar(
	cfg *config.Config,
	wsServer *wsServer.WsServer,
) *chatWSHandlers.ChatWsRegistrar {
	wire.Build(
		ClientsSet,
		ChatSet,
		chatWSHandlers.NewChatWSHandlers,
		chatWSHandlers.NewChatWsRegistrar,
	)
	return &chatWSHandlers.ChatWsRegistrar{}
}

func InitializeChatHttpServer(cfg *config.Config, wsServer *wsServer.WsServer) *httpServer.HttpServer {
	wire.Build(
		ClientsSet,
		ChatSet,
		AllHttpRegistrars,
		httpServer.New,
		authmiddleware.NewAuthMiddleware,
		corsmiddleware.NewCORSMiddleware,
	)
	return &httpServer.HttpServer{}
}

func InitializeChatGrpcServer(cfg *config.Config, wsServer *wsServer.WsServer) *grpcServer.GrpcServer {
	wire.Build(
		ClientsSet,
		ChatSet,
		AllGrpcRegistrars,
		grpcServer.New,
	)
	return &grpcServer.GrpcServer{}
}
