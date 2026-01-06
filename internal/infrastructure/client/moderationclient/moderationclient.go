package moderationclient

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/F0urward/proftwist-backend/config"
)

func NewModerationClient(cfg *config.Config) ModerationServiceClient {
	connStr := fmt.Sprintf("%s%s", cfg.ServiceHosts.Moderation, cfg.Service.GRPC.Port)
	conn, err := grpc.NewClient(connStr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	client := NewModerationServiceClient(conn)

	return client
}
