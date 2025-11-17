package main

import (
	"github.com/F0urward/proftwist-backend/config"
	roadmapinfoWire "github.com/F0urward/proftwist-backend/internal/wire/roadmapinfo"
)

func main() {
	cfg := config.New()

	httpServer := roadmapinfoWire.InitializeRoadmapInfoHttpServer(cfg)

	grpcServer := roadmapinfoWire.InitializeRoadmapInfoGrpcServer(cfg)

	go httpServer.Run()

	grpcServer.Run()
}
