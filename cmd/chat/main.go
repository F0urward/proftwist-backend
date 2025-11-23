package main

import (
	"context"

	"github.com/F0urward/proftwist-backend/config"
	chatWire "github.com/F0urward/proftwist-backend/internal/wire/chat"
)

func main() {
	cfg := config.New()

	wsServer := chatWire.InitializeChatWsServer(cfg)

	httpServer := chatWire.InitializeChatHttpServer(cfg, wsServer)

	grpcServer := chatWire.InitializeChatGrpcServer(cfg, wsServer)

	notificationWorker := chatWire.InitializeNotificationWorker(cfg, wsServer)

	notificationWorker.Start(context.Background())

	go wsServer.Run()

	go httpServer.Run()

	grpcServer.Run()
}
