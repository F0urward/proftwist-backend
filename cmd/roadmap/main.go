package main

import (
	"github.com/F0urward/proftwist-backend/config"
	roadmapWire "github.com/F0urward/proftwist-backend/internal/wire/roadmap"
	"github.com/F0urward/proftwist-backend/pkg/logger/logrus"
)

func main() {
	cfg := config.New()
	log := logrus.NewLogrusLogger()

	httpServer := roadmapWire.InitializeRoadmapHttpServer(cfg, log)

	grpcServer := roadmapWire.InitializeRoadmapGrpcServer(cfg, log)

	go httpServer.Run()

	grpcServer.Run()
}
