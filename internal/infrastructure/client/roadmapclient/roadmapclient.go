package roadmapclient

import (
	"fmt"
	"log"

	"google.golang.org/grpc"

	"github.com/F0urward/proftwist-backend/config"
	"google.golang.org/grpc/credentials/insecure"
)

func NewRoadmapClient(cfg *config.Config) RoadmapServiceClient {
	connStr := fmt.Sprintf("%s%s", cfg.Service.Host, cfg.Service.GRPC.Port)
	conn, err := grpc.NewClient(connStr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	client := NewRoadmapServiceClient(conn)

	return client
}
