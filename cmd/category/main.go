package main

import (
	"net/http"

	"github.com/F0urward/proftwist-backend/config"
	categoryWire "github.com/F0urward/proftwist-backend/internal/wire/category"
	"github.com/F0urward/proftwist-backend/pkg/logger/logrus"
	"github.com/gorilla/mux"
)

func main() {
	cfg := config.New()
	log := logrus.NewLogrusLogger()
	metrics := categoryWire.InitializeMetrics()

	go func() {
		mux := mux.NewRouter()
		mux.Handle("/metrics", metrics.Handler())
		log.Fatal(http.ListenAndServe(cfg.Metrics.Category.Port, mux))
	}()

	httpServer := categoryWire.InitializeCategoryHttpServer(cfg, log, metrics)

	httpServer.Run()
}
