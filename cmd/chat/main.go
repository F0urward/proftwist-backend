package main

import (
	"context"

	"github.com/F0urward/proftwist-backend/config"
	chatWire "github.com/F0urward/proftwist-backend/internal/wire/chat"
	"github.com/F0urward/proftwist-backend/pkg/logger/logrus"
)

func main() {
	cfg := config.New()
	log := logrus.NewLogrusLogger()

	wsServer := chatWire.InitializeChatWsServer(cfg, log)

	httpServer := chatWire.InitializeChatHttpServer(cfg, wsServer, log)

	grpcServer := chatWire.InitializeChatGrpcServer(cfg, wsServer, log)

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
