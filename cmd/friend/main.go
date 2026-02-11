package main

import (
	"net/http"

	"github.com/F0urward/proftwist-backend/config"
	friendWire "github.com/F0urward/proftwist-backend/internal/wire/friend"
	"github.com/F0urward/proftwist-backend/pkg/logger/logrus"
	"github.com/gorilla/mux"
)

func main() {
	cfg := config.New()
	log := logrus.NewLogrusLogger()
	metrics := friendWire.InitializeMetrics()

	go func() {
		mux := mux.NewRouter()
		mux.Handle("/metrics", metrics.Handler())
		log.Fatal(http.ListenAndServe(cfg.Metrics.Friend.Port, mux))
	}()

	httpServer := friendWire.InitializeFriendHttpServer(cfg, log, metrics)

	grpcServer := friendWire.InitializeFriendGrpcServer(cfg, log, metrics)

	go httpServer.Run()

	grpcServer.Run()
}
