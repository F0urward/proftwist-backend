package main

import (
	"net/http"

	"github.com/F0urward/proftwist-backend/config"
	authWire "github.com/F0urward/proftwist-backend/internal/wire/auth"
	"github.com/F0urward/proftwist-backend/pkg/logger/logrus"
	"github.com/gorilla/mux"
)

func main() {
	cfg := config.New()
	log := logrus.NewLogrusLogger()
	metrics := authWire.InitializeMetrics()

	go func() {
		mux := mux.NewRouter()
		mux.Handle("/metrics", metrics.Handler())
		log.Fatal(http.ListenAndServe(cfg.Metrics.Auth.Port, mux))
	}()

	httpServer := authWire.InitializeAuthHttpServer(cfg, log, metrics)

	grpcServer := authWire.InitializeAuthGrpcServer(cfg, log, metrics)

	go httpServer.Run()

	grpcServer.Run()
}
