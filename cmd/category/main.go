package main

import (
	"github.com/F0urward/proftwist-backend/config"
	categoryWire "github.com/F0urward/proftwist-backend/internal/wire/category"
	"github.com/F0urward/proftwist-backend/pkg/logger/logrus"
)

func main() {
	cfg := config.New()
	log := logrus.NewLogrusLogger()

	httpServer := categoryWire.InitializeCategoryHttpServer(cfg, log)

	httpServer.Run()
}
