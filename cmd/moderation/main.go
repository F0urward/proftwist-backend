package main

import (
	"github.com/F0urward/proftwist-backend/config"
	moderationWire "github.com/F0urward/proftwist-backend/internal/wire/moderation"
	"github.com/F0urward/proftwist-backend/pkg/logger/logrus"
)

func main() {
	cfg := config.New()
	log := logrus.NewLogrusLogger()

	grpcServer := moderationWire.InitializeModerationGrpcServer(cfg, log)

	grpcServer.Run()
}
