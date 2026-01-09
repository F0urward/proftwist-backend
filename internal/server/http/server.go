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
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logging"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
)

const (
	ctxTimeout = 5
)

type HttpRegistrar interface {
	RegisterRoutes(s *HttpServer)
}

type HttpServer struct {
	CFG            *config.Config
	MUX            *mux.Router
	Server         *http.Server
	AuthMiddleware *authmiddleware.AuthMiddleware
	Registrars     []HttpRegistrar
}

func (s *HttpServer) RegisterHandlers() {
	for _, registrar := range s.Registrars {
		registrar.RegisterRoutes(s)
	}
}

func New(
	cfg *config.Config,
	authMiddleware *authmiddleware.AuthMiddleware,
	corsMiddleware *corsmiddleware.CORSMiddleware,
	loggingMiddleware *logging.LoggingMiddleware,
	registrars ...HttpRegistrar,
) *HttpServer {
	mux := mux.NewRouter()

	corsedMux := corsMiddleware.CORSMiddleware(mux)
	loggedCorsedMux := loggingMiddleware.LoggingMiddleware(corsedMux)

	return &HttpServer{
		CFG: cfg,
		MUX: mux,
		Server: &http.Server{
			Addr:    cfg.Service.HTTP.Port,
			Handler: loggedCorsedMux,
		},
		AuthMiddleware: authMiddleware,
		Registrars:     registrars,
	}
}

func (s *HttpServer) Run() {
	const op = "HttpServer.Run"
	logger := ctxutil.GetLogger(context.Background()).WithField("op", op)

	s.RegisterHandlers()

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
