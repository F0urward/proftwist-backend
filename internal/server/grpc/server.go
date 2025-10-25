package grpc

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	roadmapService "github.com/F0urward/proftwist-backend/services/roadmap/delivery/grpc"
	roadmapInfoService "github.com/F0urward/proftwist-backend/services/roadmapinfo/delivery/grpc"
)

type GrpcServer struct {
	CFG               *config.Config
	Server            *grpc.Server
	RoadmapServer     *roadmapService.RoadmapServer
	RoadmapInfoServer *roadmapInfoService.RoadmapInfoServer
}

func New(
	cfg *config.Config,
	roadmapServer *roadmapService.RoadmapServer,
	roadmapInfoServer *roadmapInfoService.RoadmapInfoServer,
) *GrpcServer {
	return &GrpcServer{
		CFG:               cfg,
		Server:            grpc.NewServer(),
		RoadmapServer:     roadmapServer,
		RoadmapInfoServer: roadmapInfoServer,
	}
}

func (s *GrpcServer) Run() {
	const op = "GrpcServer.Run"
	ctx := context.Background()
	logger := logctx.GetLogger(ctx).WithField("op", op)

	l, err := net.Listen("tcp", s.CFG.Service.GRPC.Port)
	if err != nil {
		logger.WithError(err).Fatal("Error Listen in grpc server")
	}
	defer func() {
		if err := l.Close(); err != nil {
			logger.WithError(err).Error("Error closing listener")
		}
	}()

	s.RegisterServices()

	go func() {
		logger.Infof("Starting grpc server on %s", s.CFG.Service.GRPC.Port)
		if err := s.Server.Serve(l); err != nil && err != grpc.ErrServerStopped {
			logger.WithError(err).Fatal("Error Serve in grpc server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("GRPC server graceful shutdown")
	s.Server.GracefulStop()
}
