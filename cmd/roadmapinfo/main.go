package main

import (
	"github.com/F0urward/proftwist-backend/config"
	roadmapinfoWire "github.com/F0urward/proftwist-backend/internal/wire/roadmapinfo"
	"github.com/F0urward/proftwist-backend/pkg/logger/logrus"
)

func main() {
	cfg := config.New()
	log := logrus.NewLogrusLogger()

	httpServer := roadmapinfoWire.InitializeRoadmapInfoHttpServer(cfg, log)

	grpcServer := roadmapinfoWire.InitializeRoadmapInfoGrpcServer(cfg, log)

	go httpServer.Run()

	grpcServer.Run()
}
