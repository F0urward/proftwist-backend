package roadmapclient

import (
	context "context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
)

func NewRoadmapClient(cfg *config.Config) RoadmapServiceClient {
	const op = "NewRoadmapClient"
	logger := ctxutil.GetLogger(context.Background()).WithField("op", op)

	connStr := fmt.Sprintf("%s%s", cfg.ServiceHosts.Roadmap, cfg.Service.GRPC.Port)
	conn, err := grpc.NewClient(connStr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("failed to connect: %v", err)
	}

	client := NewRoadmapServiceClient(conn)

	return client
}
