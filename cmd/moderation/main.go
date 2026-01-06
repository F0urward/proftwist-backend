package main

import (
	"github.com/F0urward/proftwist-backend/config"
	moderationWire "github.com/F0urward/proftwist-backend/internal/wire/moderation"
)

func main() {
	cfg := config.New()

	grpcServer := moderationWire.InitializeModerationGrpcServer(cfg)

	grpcServer.Run()
}
