package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/F0urward/proftwist-backend/internal/server/websocket"

	"github.com/gorilla/mux"

	"github.com/F0urward/proftwist-backend/config"
	authmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/auth"
	corsmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/cors"
	"github.com/F0urward/proftwist-backend/services/auth"
	"github.com/F0urward/proftwist-backend/services/chat"
	"github.com/F0urward/proftwist-backend/services/roadmap"
	"github.com/F0urward/proftwist-backend/services/roadmapinfo"
)

const (
	ctxTimeout = 5
)

type HttpServer struct {
	CFG            *config.Config
	MUX            *mux.Router
	Server         *http.Server
	RoadmapInfoH   roadmapinfo.Handlers
	RoadmapH       roadmap.Handlers
	AuthH          auth.Handlers
	ChatH          chat.Handlers
	ChatWSH        chat.WSHandlers
	AuthMiddleware *authmiddleware.AuthMiddleware
	CORSMiddleware *corsmiddleware.CORSMiddleware
	WebSocketH     *websocket.WebSocketHandler
	WSServer       *websocket.Server
}

func New(
	cfg *config.Config,
	roadmapInfoH roadmapinfo.Handlers,
	roadmapH roadmap.Handlers,
	authH auth.Handlers,
	authMiddleware *authmiddleware.AuthMiddleware,
	chatHandler chat.Handlers,
	wsHandler *websocket.WebSocketHandler,
	wsServer *websocket.Server,
	chatwsHandler chat.WSHandlers,
	corsMiddleware *corsmiddleware.CORSMiddleware,
) *HttpServer {
	mux := mux.NewRouter()

	corsedMux := corsMiddleware.CORSMiddleware(mux)
	return &HttpServer{
		CFG:        cfg,
		MUX:        mux,
		WebSocketH: wsHandler,
		WSServer:   wsServer,
		Server: &http.Server{
			Addr:    cfg.Service.HTTP.Port,
			Handler: corsedMux,
		},
		RoadmapInfoH:   roadmapInfoH,
		RoadmapH:       roadmapH,
		AuthH:          authH,
		ChatH:          chatHandler,
		ChatWSH:        chatwsHandler,
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
