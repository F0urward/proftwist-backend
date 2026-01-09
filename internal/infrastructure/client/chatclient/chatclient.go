package chatclient

import (
	context "context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
)

func NewChatClient(cfg *config.Config) ChatServiceClient {
	const op = "NewChatClient"
	logger := ctxutil.GetLogger(context.Background()).WithField("op", op)

	connStr := fmt.Sprintf("%s%s", cfg.ServiceHosts.Chat, cfg.Service.GRPC.Port)
	conn, err := grpc.NewClient(connStr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("failed to connect: %v", err)
	}

	client := NewChatServiceClient(conn)

	return client
}
