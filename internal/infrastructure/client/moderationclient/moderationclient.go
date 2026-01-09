package moderationclient

import (
	context "context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
)

func NewModerationClient(cfg *config.Config) ModerationServiceClient {
	const op = "NewModerationClient"
	logger := ctxutil.GetLogger(context.Background()).WithField("op", op)

	connStr := fmt.Sprintf("%s%s", cfg.ServiceHosts.Moderation, cfg.Service.GRPC.Port)
	conn, err := grpc.NewClient(connStr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("failed to connect: %v", err)
	}

	client := NewModerationServiceClient(conn)

	return client
}
