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
	"github.com/F0urward/proftwist-backend/internal/worker"
)

func InitializeChatWsServer(cfg *config.Config) *wsServer.WsServer {
	wire.Build(
		ClientsSet,
		ChatSet,
		ProvideNotificationPublisher,
		ProvideBotPublisher,
		AllWsRegistrars,
		wsServer.New,
	)
	return &wsServer.WsServer{}
}

func InitializeChatHttpServer(cfg *config.Config, wsServer *wsServer.WsServer) *httpServer.HttpServer {
	wire.Build(
		ClientsSet,
		ChatSet,
		ProvideNotificationPublisher,
		ProvideBotPublisher,
		WsSet,
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
		ProvideNotificationPublisher,
		ProvideBotPublisher,
		AllGrpcRegistrars,
		grpcServer.New,
	)
	return &grpcServer.GrpcServer{}
}

func InitializeNotificationWorker(cfg *config.Config, wsServer *wsServer.WsServer) *worker.NotificationWorker {
	wire.Build(
		NotificationSet,
		ProvideNotificationConsumerConfig,
		BrokerSet,
		worker.NewNotificationWorker,
	)
	return &worker.NotificationWorker{}
}

func InitializeBotWorker(cfg *config.Config) *worker.BotWorker {
	wire.Build(
		ClientsSet,
		BotSet,
		ProvideBotConsumerConfig,
		BrokerSet,
		worker.NewBotWorker,
	)
	return &worker.BotWorker{}
}
