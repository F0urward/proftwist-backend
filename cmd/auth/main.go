package main

import (
	"github.com/F0urward/proftwist-backend/config"
	authWire "github.com/F0urward/proftwist-backend/internal/wire/auth"
	"github.com/F0urward/proftwist-backend/pkg/logger/logrus"
)

func main() {
	cfg := config.New()
	log := logrus.NewLogrusLogger()

	httpServer := authWire.InitializeAuthHttpServer(cfg, log)

	grpcServer := authWire.InitializeAuthGrpcServer(cfg, log)

	go httpServer.Run()

	grpcServer.Run()
}
