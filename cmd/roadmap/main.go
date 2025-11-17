package main

import (
	"github.com/F0urward/proftwist-backend/config"
	roadmapWire "github.com/F0urward/proftwist-backend/internal/wire/roadmap"
)

func main() {
	cfg := config.New()

	httpServer := roadmapWire.InitializeRoadmapHttpServer(cfg)

	grpcServer := roadmapWire.InitializeRoadmapGrpcServer(cfg)

	go httpServer.Run()

	grpcServer.Run()
}
