package http

import (
	"context"
	"github.com/F0urward/proftwist-backend/internal/server/websocket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/F0urward/proftwist-backend/config"
	authmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/auth"
	corsmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/cors"
	"github.com/F0urward/proftwist-backend/services/auth"
	chatdelivery "github.com/F0urward/proftwist-backend/services/chat/delivery/http"
	"github.com/F0urward/proftwist-backend/services/roadmap"
	"github.com/F0urward/proftwist-backend/services/roadmapinfo"
)

const (
	ctxTimeout = 5
)

type HttpServer struct {
	CFG                  *config.Config
	MUX                  *mux.Router
	Server               *http.Server
	RoadmapInfoH         roadmapinfo.Handlers
	RoadmapH             roadmap.Handlers
	AuthH                auth.Handlers
	ChatH                *chatdelivery.ChatHandler
	AuthMiddleware       *authmiddleware.AuthMiddleware
	CORSMiddleware       *corsmiddleware.CORSMiddleware
	WebSocketH           *chatdelivery.WebSocketHandler
	WSServer             *websocket.Server
	WebSocketIntegration *chatdelivery.WebSocketIntegration
}

func New(
	cfg *config.Config,
	roadmapInfoH roadmapinfo.Handlers,
	roadmapH roadmap.Handlers,
	authH auth.Handlers,
	authMiddleware *authmiddleware.AuthMiddleware,
	chatHandler *chatdelivery.ChatHandler,
	wsHandler *chatdelivery.WebSocketHandler,
	wsServer *websocket.Server,
	wsIntegration *chatdelivery.WebSocketIntegration,
	corsMiddleware *corsmiddleware.CORSMiddleware,
) *HttpServer {
	mux := mux.NewRouter()

	corsedMux := corsMiddleware.CORSMiddleware(mux)
	return &HttpServer{
		CFG:                  cfg,
		MUX:                  mux,
		WebSocketH:           wsHandler,
		WSServer:             wsServer,
		WebSocketIntegration: wsIntegration,
		Server: &http.Server{
			Addr:    cfg.Service.HTTP.Port,
			Handler: corsedMux,
		},
		RoadmapInfoH:   roadmapInfoH,
		RoadmapH:       roadmapH,
		AuthH:          authH,
		ChatH:          chatHandler,
		AuthMiddleware: authMiddleware,
	}
}

func (s *HttpServer) Run() {
	s.MapHandlers()

	s.WSServer.EnableDebugLogging()

	go s.WSServer.Run()

	go func() {
		log.Printf("Starting http server on %s", s.CFG.Service.HTTP.Port)
		if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error ListenAndServe in http server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	log.Println("Http server graceful shutdown")
	if err := s.Server.Shutdown(ctx); err != nil {
		log.Fatalf("Http server shutdown failed: %v", err)
	}
}
