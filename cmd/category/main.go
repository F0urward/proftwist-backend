package main

import (
	"github.com/F0urward/proftwist-backend/config"
	categoryWire "github.com/F0urward/proftwist-backend/internal/wire/category"
)

func main() {
	cfg := config.New()

	httpServer := categoryWire.InitializeCategoryHttpServer(cfg)

	httpServer.Run()
}
