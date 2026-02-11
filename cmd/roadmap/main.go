package main

import (
	"net/http"

	"github.com/F0urward/proftwist-backend/config"
	roadmapWire "github.com/F0urward/proftwist-backend/internal/wire/roadmap"
	"github.com/F0urward/proftwist-backend/pkg/logger/logrus"
	"github.com/gorilla/mux"
)

func main() {
	cfg := config.New()
	log := logrus.NewLogrusLogger()
	metrics := roadmapWire.InitializeMetrics()

	go func() {
		mux := mux.NewRouter()
		mux.Handle("/metrics", metrics.Handler())
		log.Fatal(http.ListenAndServe(cfg.Metrics.Roadmap.Port, mux))
	}()

	httpServer := roadmapWire.InitializeRoadmapHttpServer(cfg, log, metrics)

	grpcServer := roadmapWire.InitializeRoadmapGrpcServer(cfg, log, metrics)

	go httpServer.Run()

	grpcServer.Run()
}
