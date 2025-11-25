package main

import (
	"github.com/F0urward/proftwist-backend/config"
	materialWire "github.com/F0urward/proftwist-backend/internal/wire/material"
)

func main() {
	cfg := config.New()

	httpServer := materialWire.InitializeMaterialHttpServer(cfg)

	httpServer.Run()
}
