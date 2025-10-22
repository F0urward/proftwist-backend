package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/F0urward/proftwist-backend/config"
	authmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/auth"
	corsmiddleware "github.com/F0urward/proftwist-backend/internal/server/middleware/cors"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/auth"
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
	AuthMiddleware *authmiddleware.AuthMiddleware
}

func New(
	cfg *config.Config,
	roadmapInfoH roadmapinfo.Handlers,
	roadmapH roadmap.Handlers,
	authH auth.Handlers,
	authMiddleware *authmiddleware.AuthMiddleware,
	corsMiddleware *corsmiddleware.CORSMiddleware,
) *HttpServer {
	mux := mux.NewRouter()
	corsedMux := corsMiddleware.CORSMiddleware(mux)
	return &HttpServer{
		CFG: cfg,
		MUX: mux,
		Server: &http.Server{
			Addr:    cfg.Service.HTTP.Port,
			Handler: corsedMux,
		},
		RoadmapInfoH:   roadmapInfoH,
		RoadmapH:       roadmapH,
		AuthH:          authH,
		AuthMiddleware: authMiddleware,
	}
}

func (s *HttpServer) Run() {
	const op = "HttpServer.Run"
	logger := logctx.GetLogger(context.Background()).WithField("op", op)

	s.MapHandlers()

	go func() {
		logger.Infof("Starting http server on %s", s.CFG.Service.HTTP.Port)
		if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Error("Error ListenAndServe in http server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	logger.Info("Http server graceful shutdown")
	if err := s.Server.Shutdown(ctx); err != nil {
		logger.WithError(err).Fatal("Http server shutdown failed")
	}
}
