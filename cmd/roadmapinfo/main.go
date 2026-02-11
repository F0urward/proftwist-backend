package main

import (
	"net/http"

	"github.com/F0urward/proftwist-backend/config"
	roadmapinfoWire "github.com/F0urward/proftwist-backend/internal/wire/roadmapinfo"
	"github.com/F0urward/proftwist-backend/pkg/logger/logrus"
	"github.com/gorilla/mux"
)

func main() {
	cfg := config.New()
	log := logrus.NewLogrusLogger()
	metrics := roadmapinfoWire.InitializeMetrics()

	go func() {
		mux := mux.NewRouter()
		mux.Handle("/metrics", metrics.Handler())
		log.Fatal(http.ListenAndServe(cfg.Metrics.Roadmapinfo.Port, mux))
	}()

	httpServer := roadmapinfoWire.InitializeRoadmapInfoHttpServer(cfg, log, metrics)

	grpcServer := roadmapinfoWire.InitializeRoadmapInfoGrpcServer(cfg, log, metrics)

	go httpServer.Run()

	grpcServer.Run()
}
