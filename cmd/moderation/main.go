package main

import (
	"net/http"

	"github.com/F0urward/proftwist-backend/config"
	moderationWire "github.com/F0urward/proftwist-backend/internal/wire/moderation"
	"github.com/F0urward/proftwist-backend/pkg/logger/logrus"
	"github.com/gorilla/mux"
)

func main() {
	cfg := config.New()
	log := logrus.NewLogrusLogger()
	metrics := moderationWire.InitializeMetrics()

	go func() {
		mux := mux.NewRouter()
		mux.Handle("/metrics", metrics.Handler())
		log.Fatal(http.ListenAndServe(cfg.Metrics.Moderation.Port, mux))
	}()

	grpcServer := moderationWire.InitializeModerationGrpcServer(cfg, log, metrics)

	grpcServer.Run()
}
