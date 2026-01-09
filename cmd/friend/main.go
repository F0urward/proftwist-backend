package main

import (
	"github.com/F0urward/proftwist-backend/config"
	friendWire "github.com/F0urward/proftwist-backend/internal/wire/friend"
	"github.com/F0urward/proftwist-backend/pkg/logger/logrus"
)

func main() {
	cfg := config.New()
	log := logrus.NewLogrusLogger()

	httpServer := friendWire.InitializeFriendHttpServer(cfg, log)

	grpcServer := friendWire.InitializeFriendGrpcServer(cfg, log)

	go httpServer.Run()

	grpcServer.Run()
}
