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

	for i := 0; i < cfg.Workers.Notification.Count; i++ {
		notificationWorker := chatWire.InitializeNotificationWorker(cfg, wsServer)
		go notificationWorker.Start(context.Background())
	}

	for i := 0; i < cfg.Workers.Bot.Count; i++ {
		botWorker := chatWire.InitializeBotWorker(cfg)
		go botWorker.Start(context.Background())
	}

	go wsServer.Run()

	go httpServer.Run()

	grpcServer.Run()
}
