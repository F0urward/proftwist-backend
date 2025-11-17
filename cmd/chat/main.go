package main

import (
	"github.com/F0urward/proftwist-backend/config"
	chatWire "github.com/F0urward/proftwist-backend/internal/wire/chat"
)

func main() {
	cfg := config.New()

	wsServer := chatWire.InitializeChatWsServer(cfg)
	chatWsRegistrar := chatWire.IntitializeChatWsRegistrar(cfg, wsServer)
	chatWsRegistrar.RegisterHandlers(wsServer)

	httpServer := chatWire.InitializeChatHttpServer(cfg, wsServer)

	grpcServer := chatWire.InitializeChatGrpcServer(cfg, wsServer)

	go wsServer.Run()

	go httpServer.Run()

	grpcServer.Run()
}
