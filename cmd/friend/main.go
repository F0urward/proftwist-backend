package main

import (
	"github.com/F0urward/proftwist-backend/config"
	friendWire "github.com/F0urward/proftwist-backend/internal/wire/friend"
)

func main() {
	cfg := config.New()

	httpServer := friendWire.InitializeFriendHttpServer(cfg)

	httpServer.Run()
}
