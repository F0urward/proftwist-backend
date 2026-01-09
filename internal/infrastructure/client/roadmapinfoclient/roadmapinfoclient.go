package roadmapinfoclient

import (
	context "context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
)

func NewRoadmapInfoClient(cfg *config.Config) RoadmapInfoServiceClient {
	const op = "NewRoadmapInfoClient"
	logger := ctxutil.GetLogger(context.Background()).WithField("op", op)

	connStr := fmt.Sprintf("%s%s", cfg.ServiceHosts.RoadmapInfo, cfg.Service.GRPC.Port)
	conn, err := grpc.NewClient(connStr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("failed to connect: %v", err)
	}

	client := NewRoadmapInfoServiceClient(conn)

	return client
}
