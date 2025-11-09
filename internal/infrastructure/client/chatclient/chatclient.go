package chatclient

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/F0urward/proftwist-backend/config"
)

func NewChatClient(cfg *config.Config) ChatServiceClient {
	connStr := fmt.Sprintf("%s%s", cfg.Service.Host, cfg.Service.GRPC.Port)
	conn, err := grpc.NewClient(connStr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	client := NewChatServiceClient(conn)

	return client
}
