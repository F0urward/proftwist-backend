package friendclient

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/F0urward/proftwist-backend/config"
)

func NewFriendClient(cfg *config.Config) FriendServiceClient {
	connStr := fmt.Sprintf("%s%s", cfg.ServiceHosts.Friend, cfg.Service.GRPC.Port)
	conn, err := grpc.NewClient(connStr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	client := NewFriendServiceClient(conn)

	return client
}
